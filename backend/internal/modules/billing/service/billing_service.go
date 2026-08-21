package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"queue-management-tenant/backend/internal/modules/billing/dto"
	"queue-management-tenant/backend/internal/modules/billing/entity"
	billingRepo "queue-management-tenant/backend/internal/modules/billing/repository"
	orgRepo "queue-management-tenant/backend/internal/modules/organization/repository"
	"queue-management-tenant/backend/pkg/midtrans"
)

type BillingService struct {
	repo           *billingRepo.BillingRepository
	orgRepo        *orgRepo.OrganizationRepository
	midtransClient *midtrans.MidtransClient
}

func NewBillingService(repo *billingRepo.BillingRepository, orgRepo *orgRepo.OrganizationRepository) *BillingService {
	return &BillingService{
		repo:           repo,
		orgRepo:        orgRepo,
		midtransClient: midtrans.NewMidtransClient(),
	}
}

func (s *BillingService) IncrementTicketCount(ctx context.Context, orgID int64) error {
	period := time.Now().Format("2006-01")
	return s.repo.IncrementTicketMeter(ctx, orgID, period)
}

func (s *BillingService) GetCurrentUsage(ctx context.Context, orgID int64) (*entity.UsageMeter, error) {
	period := time.Now().Format("2006-01")
	return s.repo.GetUsageMeter(ctx, orgID, period)
}

func (s *BillingService) ListTenantInvoices(ctx context.Context, orgID int64) ([]entity.Invoice, error) {
	return s.repo.ListInvoicesByOrg(ctx, orgID)
}

func (s *BillingService) ListAllInvoices(ctx context.Context) ([]entity.Invoice, error) {
	return s.repo.ListAllInvoices(ctx)
}

func (s *BillingService) CreateSnapToken(ctx context.Context, orgID int64, invoiceID int64) (*dto.CreateSnapTokenResponse, error) {
	inv, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	if inv.OrganizationID != orgID {
		return nil, fmt.Errorf("unauthorized invoice access")
	}

	if inv.Status == "PAID" {
		return nil, fmt.Errorf("invoice is already paid")
	}

	// Check if there is an active pending payment with snap token
	if inv.LatestPayment != nil && inv.LatestPayment.Status == "PENDING" && inv.LatestPayment.SnapToken != "" {
		return &dto.CreateSnapTokenResponse{
			PaymentNumber:   inv.LatestPayment.PaymentNumber,
			SnapToken:       inv.LatestPayment.SnapToken,
			SnapRedirectURL: inv.LatestPayment.SnapRedirectURL,
		}, nil
	}

	// Create new Payment record & Midtrans Snap Token
	paymentNumber := fmt.Sprintf("PAY-%s-%d-%d", inv.BillingPeriod, inv.ID, time.Now().Unix())
	amount := int64(inv.TotalAmount)
	description := fmt.Sprintf("QMS SaaS Metered Bill %s (Invoice #%s)", inv.BillingPeriod, inv.InvoiceNumber)

	snapResp, err := s.midtransClient.CreateSnapToken(paymentNumber, amount, inv.OrgName, "billing@tenant.com", description)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Midtrans Snap token: %w", err)
	}

	payment := &entity.Payment{
		PaymentNumber:   paymentNumber,
		InvoiceID:       inv.ID,
		OrganizationID: orgID,
		Amount:          inv.TotalAmount,
		Status:          "PENDING",
		SnapToken:       snapResp.Token,
		SnapRedirectURL: snapResp.RedirectURL,
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to record payment transaction: %w", err)
	}

	return &dto.CreateSnapTokenResponse{
		PaymentNumber:   paymentNumber,
		SnapToken:       snapResp.Token,
		SnapRedirectURL: snapResp.RedirectURL,
	}, nil
}

func (s *BillingService) ProcessMidtransWebhook(ctx context.Context, payload dto.MidtransWebhookPayload) error {
	rawBytes, _ := json.Marshal(payload)
	_ = s.repo.SaveWebhookLog(ctx, payload.OrderID, payload.TransactionStatus, string(rawBytes))

	status := "PENDING"
	if payload.TransactionStatus == "settlement" || payload.TransactionStatus == "capture" {
		status = "SETTLEMENT"
	} else if payload.TransactionStatus == "deny" || payload.TransactionStatus == "cancel" || payload.TransactionStatus == "expire" {
		status = "CANCELLED"
	}

	return s.repo.UpdatePaymentStatusTx(ctx, payload.OrderID, status, payload.PaymentType, payload.PaymentType, string(rawBytes))
}

func (s *BillingService) GetSuperadminStats(ctx context.Context) (*dto.SuperadminBillingStatsResponse, error) {
	return s.repo.GetSuperadminStats(ctx)
}
