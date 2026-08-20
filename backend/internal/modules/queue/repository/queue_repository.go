package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"queue-management-tenant/backend/internal/modules/queue/entity"
)

type QueueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) *QueueRepository {
	return &QueueRepository{db: db}
}

// IssueTicketTx generates sequence number atomically using FOR UPDATE pessimistic lock
func (r *QueueRepository) IssueTicketTx(ctx context.Context, orgID, branchID, serviceID int64, customerID sql.NullInt64, priority string, prefix string, avgServiceTimeSec int) (*entity.QueueTicket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	today := time.Now().UTC().Format("2006-01-02")

	upsertQuery := `
		INSERT INTO queue_sequences (organization_id, branch_id, service_id, sequence_date, last_number)
		VALUES ($1, $2, $3, $4, 0)
		ON CONFLICT (branch_id, service_id, sequence_date) DO NOTHING
	`
	if _, err := tx.ExecContext(ctx, upsertQuery, orgID, branchID, serviceID, today); err != nil {
		return nil, fmt.Errorf("failed sequence upsert: %w", err)
	}

	var currentSeq int
	lockQuery := `
		SELECT last_number FROM queue_sequences
		WHERE branch_id = $1 AND service_id = $2 AND sequence_date = $3
		FOR UPDATE
	`
	if err := tx.QueryRowContext(ctx, lockQuery, branchID, serviceID, today).Scan(&currentSeq); err != nil {
		return nil, fmt.Errorf("failed to lock sequence row: %w", err)
	}

	nextSeq := currentSeq + 1

	updateSeqQuery := `
		UPDATE queue_sequences SET last_number = $1, updated_at = CURRENT_TIMESTAMP
		WHERE branch_id = $2 AND service_id = $3 AND sequence_date = $4
	`
	if _, err := tx.ExecContext(ctx, updateSeqQuery, nextSeq, branchID, serviceID, today); err != nil {
		return nil, fmt.Errorf("failed to update sequence: %w", err)
	}

	ticketNum := fmt.Sprintf("%s%03d", prefix, nextSeq)

	var peopleAhead int
	countAheadQuery := `
		SELECT COUNT(*) FROM queue_tickets
		WHERE branch_id = $1 AND service_id = $2 AND status = 'WAITING' AND queue_date = $3
	`
	_ = tx.QueryRowContext(ctx, countAheadQuery, branchID, serviceID, today).Scan(&peopleAhead)

	estimatedWaitSec := peopleAhead * avgServiceTimeSec
	if estimatedWaitSec < avgServiceTimeSec && peopleAhead == 0 {
		estimatedWaitSec = avgServiceTimeSec / 2
	}

	var cID interface{} = nil
	if customerID.Valid {
		cID = customerID.Int64
	}

	ticket := &entity.QueueTicket{}
	insertTicketQuery := `
		INSERT INTO queue_tickets (
			organization_id, branch_id, service_id, customer_id, ticket_number,
			sequence_number, queue_date, priority, status, estimated_wait_seconds
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'WAITING', $9)
		RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, estimated_wait_seconds, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, insertTicketQuery,
		orgID, branchID, serviceID, cID, ticketNum,
		nextSeq, today, priority, estimatedWaitSec,
	).Scan(
		&ticket.ID, &ticket.UUID, &ticket.OrganizationID, &ticket.BranchID, &ticket.ServiceID,
		&ticket.TicketNumber, &ticket.SequenceNumber, &ticket.QueueDate, &ticket.Priority,
		&ticket.Status, &ticket.PublicToken, &ticket.EstimatedWaitSeconds, &ticket.CreatedAt, &ticket.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert ticket: %w", err)
	}

	eventQuery := `INSERT INTO queue_ticket_events (ticket_id, from_status, to_status) VALUES ($1, NULL, 'WAITING')`
	if _, err := tx.ExecContext(ctx, eventQuery, ticket.ID); err != nil {
		return nil, fmt.Errorf("failed to log ticket event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	ticket.PeopleAhead = peopleAhead
	return ticket, nil
}

// CallNextTx selects the next waiting ticket atomically using FOR UPDATE SKIP LOCKED
func (r *QueueRepository) CallNextTx(ctx context.Context, orgID, branchID, counterID, staffID int64, serviceIDs []int64) (*entity.QueueTicket, error) {
	if len(serviceIDs) == 0 {
		return nil, sql.ErrNoRows
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	today := time.Now().UTC().Format("2006-01-02")

	query := `
		SELECT id FROM queue_tickets
		WHERE organization_id = $1 AND branch_id = $2 AND service_id = ANY($3)
		  AND status = 'WAITING' AND queue_date = $4 AND deleted_at IS NULL
		ORDER BY 
			CASE priority 
				WHEN 'EMERGENCY' THEN 1 
				WHEN 'PRIORITY' THEN 2 
				ELSE 3 
			END ASC,
			sequence_number ASC
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`
	var ticketID int64
	err = tx.QueryRowContext(ctx, query, orgID, branchID, serviceIDs, today).Scan(&ticketID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	updateQuery := `
		UPDATE queue_tickets
		SET status = 'CALLED', counter_id = $1, staff_id = $2, called_at = $3, updated_at = $3
		WHERE id = $4
		RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, called_at, created_at, updated_at
	`
	ticket := &entity.QueueTicket{}
	err = tx.QueryRowContext(ctx, updateQuery, counterID, staffID, now, ticketID).Scan(
		&ticket.ID, &ticket.UUID, &ticket.OrganizationID, &ticket.BranchID, &ticket.ServiceID,
		&ticket.TicketNumber, &ticket.SequenceNumber, &ticket.QueueDate, &ticket.Priority,
		&ticket.Status, &ticket.PublicToken, &ticket.CalledAt, &ticket.CreatedAt, &ticket.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	_, _ = tx.ExecContext(ctx, `UPDATE counters SET status = 'BUSY' WHERE id = $1`, counterID)
	_, _ = tx.ExecContext(ctx, `INSERT INTO queue_ticket_events (ticket_id, from_status, to_status, actor_id) VALUES ($1, 'WAITING', 'CALLED', $2)`, ticketID, staffID)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *QueueRepository) UpdateTicketStatus(ctx context.Context, ticketID int64, fromStatus, toStatus string, staffID int64, counterID sql.NullInt64) (*entity.QueueTicket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	var updateQuery string
	switch toStatus {
	case "SERVING":
		updateQuery = `UPDATE queue_tickets SET status = $1, serving_started_at = $2, updated_at = $2 WHERE id = $3 AND status = $4 RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, called_at, serving_started_at, created_at, updated_at`
	case "COMPLETED":
		updateQuery = `UPDATE queue_tickets SET status = $1, completed_at = $2, updated_at = $2 WHERE id = $3 AND status IN ('CALLED', 'SERVING') RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, called_at, serving_started_at, completed_at, created_at, updated_at`
	case "SKIPPED", "NO_SHOW", "CANCELLED":
		updateQuery = `UPDATE queue_tickets SET status = $1, cancelled_at = $2, updated_at = $2 WHERE id = $3 RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, called_at, serving_started_at, completed_at, cancelled_at, created_at, updated_at`
	default:
		updateQuery = `UPDATE queue_tickets SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id, uuid, organization_id, branch_id, service_id, ticket_number, sequence_number, queue_date, priority, status, public_token, created_at, updated_at`
	}

	ticket := &entity.QueueTicket{}
	err = tx.QueryRowContext(ctx, updateQuery, toStatus, now, ticketID, fromStatus).Scan(
		&ticket.ID, &ticket.UUID, &ticket.OrganizationID, &ticket.BranchID, &ticket.ServiceID,
		&ticket.TicketNumber, &ticket.SequenceNumber, &ticket.QueueDate, &ticket.Priority,
		&ticket.Status, &ticket.PublicToken, &ticket.CreatedAt, &ticket.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid ticket state transition or ticket not found: %w", err)
	}

	_, _ = tx.ExecContext(ctx, `INSERT INTO queue_ticket_events (ticket_id, from_status, to_status, actor_id) VALUES ($1, $2, $3, $4)`, ticketID, fromStatus, toStatus, staffID)

	if counterID.Valid && (toStatus == "COMPLETED" || toStatus == "SKIPPED" || toStatus == "NO_SHOW") {
		_, _ = tx.ExecContext(ctx, `UPDATE counters SET status = 'OPEN' WHERE id = $1`, counterID.Int64)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *QueueRepository) GetByPublicToken(ctx context.Context, tokenUUID string) (*entity.QueueTicket, error) {
	query := `
		SELECT t.id, t.uuid, t.organization_id, t.branch_id, b.name, t.service_id, s.name, s.prefix, t.ticket_number, t.sequence_number, t.queue_date, t.priority, t.status, t.public_token, COALESCE(c.counter_number, ''), t.called_at, t.serving_started_at, t.completed_at, t.estimated_wait_seconds, t.created_at, t.updated_at
		FROM queue_tickets t
		JOIN branches b ON t.branch_id = b.id
		JOIN services s ON t.service_id = s.id
		LEFT JOIN counters c ON t.counter_id = c.id
		WHERE t.public_token = $1 AND t.deleted_at IS NULL
	`
	t := &entity.QueueTicket{}
	err := r.db.QueryRowContext(ctx, query, tokenUUID).Scan(
		&t.ID, &t.UUID, &t.OrganizationID, &t.BranchID, &t.BranchName, &t.ServiceID, &t.ServiceName, &t.ServicePrefix,
		&t.TicketNumber, &t.SequenceNumber, &t.QueueDate, &t.Priority, &t.Status, &t.PublicToken, &t.CounterNumber,
		&t.CalledAt, &t.ServingStartedAt, &t.CompletedAt, &t.EstimatedWaitSeconds, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if t.Status == "WAITING" {
		today := t.QueueDate.Format("2006-01-02")
		_ = r.db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM queue_tickets
			WHERE branch_id = $1 AND service_id = $2 AND status = 'WAITING' AND queue_date = $3 AND sequence_number < $4
		`, t.BranchID, t.ServiceID, today, t.SequenceNumber).Scan(&t.PeopleAhead)
	}

	return t, nil
}

func (r *QueueRepository) GetByID(ctx context.Context, orgID, id int64) (*entity.QueueTicket, error) {
	query := `
		SELECT t.id, t.uuid, t.organization_id, t.branch_id, b.name, t.service_id, s.name, s.prefix, t.ticket_number, t.sequence_number, t.queue_date, t.priority, t.status, t.public_token, t.counter_id, COALESCE(c.counter_number, ''), t.staff_id, COALESCE(u.full_name, ''), t.called_at, t.serving_started_at, t.completed_at, t.estimated_wait_seconds, t.created_at, t.updated_at
		FROM queue_tickets t
		JOIN branches b ON t.branch_id = b.id
		JOIN services s ON t.service_id = s.id
		LEFT JOIN counters c ON t.counter_id = c.id
		LEFT JOIN users u ON t.staff_id = u.id
		WHERE t.id = $1 AND t.organization_id = $2 AND t.deleted_at IS NULL
	`
	t := &entity.QueueTicket{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).Scan(
		&t.ID, &t.UUID, &t.OrganizationID, &t.BranchID, &t.BranchName, &t.ServiceID, &t.ServiceName, &t.ServicePrefix,
		&t.TicketNumber, &t.SequenceNumber, &t.QueueDate, &t.Priority, &t.Status, &t.PublicToken, &t.CounterID, &t.CounterNumber, &t.StaffID, &t.StaffName,
		&t.CalledAt, &t.ServingStartedAt, &t.CompletedAt, &t.EstimatedWaitSeconds, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *QueueRepository) ListTickets(ctx context.Context, orgID, branchID int64, status string, limit int) ([]entity.QueueTicket, error) {
	today := time.Now().UTC().Format("2006-01-02")
	query := `
		SELECT t.id, t.uuid, t.organization_id, t.branch_id, b.name, t.service_id, s.name, s.prefix, t.ticket_number, t.sequence_number, t.queue_date, t.priority, t.status, t.public_token, COALESCE(c.counter_number, ''), COALESCE(u.full_name, ''), t.called_at, t.serving_started_at, t.completed_at, t.estimated_wait_seconds, t.created_at, t.updated_at
		FROM queue_tickets t
		JOIN branches b ON t.branch_id = b.id
		JOIN services s ON t.service_id = s.id
		LEFT JOIN counters c ON t.counter_id = c.id
		LEFT JOIN users u ON t.staff_id = u.id
		WHERE t.organization_id = $1 AND t.branch_id = $2 AND t.queue_date = $3 AND t.deleted_at IS NULL
	`
	args := []interface{}{orgID, branchID, today}
	if status != "" {
		query += " AND t.status = $4"
		args = append(args, status)
	}

	query += " ORDER BY t.id DESC LIMIT " + fmt.Sprintf("%d", limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []entity.QueueTicket
	for rows.Next() {
		var t entity.QueueTicket
		if err := rows.Scan(
			&t.ID, &t.UUID, &t.OrganizationID, &t.BranchID, &t.BranchName, &t.ServiceID, &t.ServiceName, &t.ServicePrefix,
			&t.TicketNumber, &t.SequenceNumber, &t.QueueDate, &t.Priority, &t.Status, &t.PublicToken, &t.CounterNumber, &t.StaffName,
			&t.CalledAt, &t.ServingStartedAt, &t.CompletedAt, &t.EstimatedWaitSeconds, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}
	return tickets, nil
}
