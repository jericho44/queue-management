package service

import (
	"context"
	"database/sql"
	"fmt"

	counterRepo "queue-management-tenant/backend/internal/modules/counter/repository"
	"queue-management-tenant/backend/internal/modules/queue/dto"
	queueEntity "queue-management-tenant/backend/internal/modules/queue/entity"
	queueRepo "queue-management-tenant/backend/internal/modules/queue/repository"
	serviceRepo "queue-management-tenant/backend/internal/modules/service/repository"
	"queue-management-tenant/backend/internal/websocket"
)

type QueueService struct {
	queueRepo   *queueRepo.QueueRepository
	serviceRepo *serviceRepo.ServiceRepository
	counterRepo *counterRepo.CounterRepository
	wsHub       *websocket.WSHub
}

func NewQueueService(queueRepo *queueRepo.QueueRepository, serviceRepo *serviceRepo.ServiceRepository, counterRepo *counterRepo.CounterRepository, wsHub *websocket.WSHub) *QueueService {
	return &QueueService{
		queueRepo:   queueRepo,
		serviceRepo: serviceRepo,
		counterRepo: counterRepo,
		wsHub:       wsHub,
	}
}

func (s *QueueService) IssueTicket(ctx context.Context, orgID int64, req dto.IssueTicketRequest) (*queueEntity.QueueTicket, error) {
	svc, err := s.serviceRepo.GetByID(ctx, orgID, req.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	priority := req.Priority
	if priority != "EMERGENCY" && priority != "PRIORITY" {
		priority = "NORMAL"
	}

	ticket, err := s.queueRepo.IssueTicketTx(ctx, orgID, req.BranchID, req.ServiceID, sql.NullInt64{Valid: false}, priority, svc.Prefix, svc.AvgServiceTimeSec)
	if err != nil {
		return nil, err
	}

	ticket.ServiceName = svc.Name
	ticket.ServicePrefix = svc.Prefix

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, req.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_CREATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) CallNext(ctx context.Context, orgID, counterID, staffID int64) (*queueEntity.QueueTicket, error) {
	counter, err := s.counterRepo.GetByID(ctx, orgID, counterID)
	if err != nil {
		return nil, fmt.Errorf("counter not found: %w", err)
	}

	if len(counter.Services) == 0 {
		return nil, fmt.Errorf("no services configured for this counter")
	}

	var serviceIDs []int64
	for _, svc := range counter.Services {
		serviceIDs = append(serviceIDs, svc.ID)
	}

	ticket, err := s.queueRepo.CallNextTx(ctx, orgID, counter.BranchID, counterID, staffID, serviceIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no waiting tickets available")
		}
		return nil, err
	}

	ticket.CounterNumber = counter.CounterNumber

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, counter.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_CALLED", ticket)
		s.wsHub.BroadcastToChannel(fmt.Sprintf("ticket:%s", ticket.PublicToken), "TICKET_STATUS_UPDATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) RecallTicket(ctx context.Context, orgID, ticketID int64) (*queueEntity.QueueTicket, error) {
	ticket, err := s.queueRepo.GetByID(ctx, orgID, ticketID)
	if err != nil {
		return nil, err
	}

	if ticket.Status != "CALLED" {
		return nil, fmt.Errorf("ticket cannot be recalled in state %s", ticket.Status)
	}

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, ticket.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_RECALLED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) StartServing(ctx context.Context, orgID, ticketID, staffID int64) (*queueEntity.QueueTicket, error) {
	ticket, err := s.queueRepo.UpdateTicketStatus(ctx, ticketID, "CALLED", "SERVING", staffID, sql.NullInt64{Valid: false})
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, ticket.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_STARTED", ticket)
		s.wsHub.BroadcastToChannel(fmt.Sprintf("ticket:%s", ticket.PublicToken), "TICKET_STATUS_UPDATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) CompleteTicket(ctx context.Context, orgID, ticketID, staffID int64) (*queueEntity.QueueTicket, error) {
	existing, err := s.queueRepo.GetByID(ctx, orgID, ticketID)
	if err != nil {
		return nil, err
	}

	ticket, err := s.queueRepo.UpdateTicketStatus(ctx, ticketID, existing.Status, "COMPLETED", staffID, existing.CounterID)
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, ticket.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_COMPLETED", ticket)
		s.wsHub.BroadcastToChannel(fmt.Sprintf("ticket:%s", ticket.PublicToken), "TICKET_STATUS_UPDATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) SkipTicket(ctx context.Context, orgID, ticketID, staffID int64) (*queueEntity.QueueTicket, error) {
	existing, err := s.queueRepo.GetByID(ctx, orgID, ticketID)
	if err != nil {
		return nil, err
	}

	ticket, err := s.queueRepo.UpdateTicketStatus(ctx, ticketID, existing.Status, "SKIPPED", staffID, existing.CounterID)
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, ticket.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_SKIPPED", ticket)
		s.wsHub.BroadcastToChannel(fmt.Sprintf("ticket:%s", ticket.PublicToken), "TICKET_STATUS_UPDATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) NoShowTicket(ctx context.Context, orgID, ticketID, staffID int64) (*queueEntity.QueueTicket, error) {
	existing, err := s.queueRepo.GetByID(ctx, orgID, ticketID)
	if err != nil {
		return nil, err
	}

	ticket, err := s.queueRepo.UpdateTicketStatus(ctx, ticketID, existing.Status, "NO_SHOW", staffID, existing.CounterID)
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("org:%d:branch:%d", orgID, ticket.BranchID)
	if s.wsHub != nil {
		s.wsHub.BroadcastToChannel(channel, "QUEUE_TICKET_NO_SHOW", ticket)
		s.wsHub.BroadcastToChannel(fmt.Sprintf("ticket:%s", ticket.PublicToken), "TICKET_STATUS_UPDATED", ticket)
	}

	return ticket, nil
}

func (s *QueueService) GetByPublicToken(ctx context.Context, tokenUUID string) (*queueEntity.QueueTicket, error) {
	return s.queueRepo.GetByPublicToken(ctx, tokenUUID)
}

func (s *QueueService) ListBranchTickets(ctx context.Context, orgID, branchID int64, status string) ([]queueEntity.QueueTicket, error) {
	return s.queueRepo.ListTickets(ctx, orgID, branchID, status, 50)
}
