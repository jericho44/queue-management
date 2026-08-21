package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"queue-management-tenant/backend/internal/modules/billing/dto"
	"queue-management-tenant/backend/internal/modules/billing/entity"
)

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) IncrementTicketMeter(ctx context.Context, orgID int64, period string) error {
	query := `
		INSERT INTO billing_usage_meters (organization_id, billing_period, ticket_count, last_ticket_at)
		VALUES ($1, $2, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (organization_id, billing_period)
		DO UPDATE SET
			ticket_count = billing_usage_meters.ticket_count + 1,
			last_ticket_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.ExecContext(ctx, query, orgID, period)
	return err
}

func (r *BillingRepository) GetUsageMeter(ctx context.Context, orgID int64, period string) (*entity.UsageMeter, error) {
	query := `
		SELECT id, organization_id, billing_period, ticket_count, last_ticket_at, created_at, updated_at
		FROM billing_usage_meters
		WHERE organization_id = $1 AND billing_period = $2
	`
	m := &entity.UsageMeter{}
	err := r.db.QueryRowContext(ctx, query, orgID, period).
		Scan(&m.ID, &m.OrganizationID, &m.BillingPeriod, &m.TicketCount, &m.LastTicketAt, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return &entity.UsageMeter{
			OrganizationID: orgID,
			BillingPeriod:  period,
			TicketCount:    0,
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (r *BillingRepository) CreateInvoiceWithItems(ctx context.Context, inv *entity.Invoice, items []entity.InvoiceItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryInv := `
		INSERT INTO invoices (invoice_number, organization_id, billing_period, subtotal, tax_amount, total_amount, status, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, uuid, issued_at, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, queryInv,
		inv.InvoiceNumber, inv.OrganizationID, inv.BillingPeriod,
		inv.Subtotal, inv.TaxAmount, inv.TotalAmount, inv.Status, inv.DueDate,
	).Scan(&inv.ID, &inv.UUID, &inv.IssuedAt, &inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert invoice header: %w", err)
	}

	queryItem := `
		INSERT INTO invoice_items (invoice_id, description, item_type, quantity, unit_price, total_price)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	for i := range items {
		items[i].InvoiceID = inv.ID
		err = tx.QueryRowContext(ctx, queryItem,
			inv.ID, items[i].Description, items[i].ItemType, items[i].Quantity, items[i].UnitPrice, items[i].TotalPrice,
		).Scan(&items[i].ID, &items[i].CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert invoice item: %w", err)
		}
	}

	inv.Items = items
	return tx.Commit()
}

func (r *BillingRepository) GetInvoiceByID(ctx context.Context, id int64) (*entity.Invoice, error) {
	queryInv := `
		SELECT i.id, i.uuid, i.invoice_number, i.organization_id, o.name, i.billing_period,
		       i.subtotal, i.tax_amount, i.total_amount, i.status, i.issued_at, i.due_date,
		       i.created_at, i.updated_at
		FROM invoices i
		JOIN organizations o ON i.organization_id = o.id
		WHERE i.id = $1 AND i.deleted_at IS NULL
	`
	inv := &entity.Invoice{}
	err := r.db.QueryRowContext(ctx, queryInv, id).
		Scan(&inv.ID, &inv.UUID, &inv.InvoiceNumber, &inv.OrganizationID, &inv.OrgName, &inv.BillingPeriod,
			&inv.Subtotal, &inv.TaxAmount, &inv.TotalAmount, &inv.Status, &inv.IssuedAt, &inv.DueDate,
			&inv.CreatedAt, &inv.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Fetch Items
	items, err := r.GetInvoiceItems(ctx, inv.ID)
	if err == nil {
		inv.Items = items
	}

	// Fetch Latest Payment
	payment, _ := r.GetLatestPaymentByInvoice(ctx, inv.ID)
	inv.LatestPayment = payment

	return inv, nil
}

func (r *BillingRepository) GetInvoiceItems(ctx context.Context, invoiceID int64) ([]entity.InvoiceItem, error) {
	query := `
		SELECT id, invoice_id, description, item_type, quantity, unit_price, total_price, created_at
		FROM invoice_items
		WHERE invoice_id = $1
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query, invoiceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.InvoiceItem
	for rows.Next() {
		var item entity.InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.ItemType, &item.Quantity, &item.UnitPrice, &item.TotalPrice, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *BillingRepository) ListInvoicesByOrg(ctx context.Context, orgID int64) ([]entity.Invoice, error) {
	query := `
		SELECT i.id, i.uuid, i.invoice_number, i.organization_id, o.name, i.billing_period,
		       i.subtotal, i.tax_amount, i.total_amount, i.status, i.issued_at, i.due_date,
		       i.created_at, i.updated_at
		FROM invoices i
		JOIN organizations o ON i.organization_id = o.id
		WHERE i.organization_id = $1 AND i.deleted_at IS NULL
		ORDER BY i.issued_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []entity.Invoice
	for rows.Next() {
		var inv entity.Invoice
		if err := rows.Scan(&inv.ID, &inv.UUID, &inv.InvoiceNumber, &inv.OrganizationID, &inv.OrgName, &inv.BillingPeriod,
			&inv.Subtotal, &inv.TaxAmount, &inv.TotalAmount, &inv.Status, &inv.IssuedAt, &inv.DueDate,
			&inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}

		items, _ := r.GetInvoiceItems(ctx, inv.ID)
		inv.Items = items

		payment, _ := r.GetLatestPaymentByInvoice(ctx, inv.ID)
		inv.LatestPayment = payment

		invoices = append(invoices, inv)
	}
	return invoices, nil
}

func (r *BillingRepository) ListAllInvoices(ctx context.Context) ([]entity.Invoice, error) {
	query := `
		SELECT i.id, i.uuid, i.invoice_number, i.organization_id, o.name, i.billing_period,
		       i.subtotal, i.tax_amount, i.total_amount, i.status, i.issued_at, i.due_date,
		       i.created_at, i.updated_at
		FROM invoices i
		JOIN organizations o ON i.organization_id = o.id
		WHERE i.deleted_at IS NULL
		ORDER BY i.issued_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []entity.Invoice
	for rows.Next() {
		var inv entity.Invoice
		if err := rows.Scan(&inv.ID, &inv.UUID, &inv.InvoiceNumber, &inv.OrganizationID, &inv.OrgName, &inv.BillingPeriod,
			&inv.Subtotal, &inv.TaxAmount, &inv.TotalAmount, &inv.Status, &inv.IssuedAt, &inv.DueDate,
			&inv.CreatedAt, &inv.UpdatedAt); err != nil {
			return nil, err
		}

		items, _ := r.GetInvoiceItems(ctx, inv.ID)
		inv.Items = items

		payment, _ := r.GetLatestPaymentByInvoice(ctx, inv.ID)
		inv.LatestPayment = payment

		invoices = append(invoices, inv)
	}
	return invoices, nil
}

// Payment Repository Operations
func (r *BillingRepository) CreatePayment(ctx context.Context, p *entity.Payment) error {
	query := `
		INSERT INTO payments (payment_number, invoice_id, organization_id, amount, payment_method, payment_channel, status, snap_token, snap_redirect_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, uuid, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		p.PaymentNumber, p.InvoiceID, p.OrganizationID, p.Amount,
		p.PaymentMethod, p.PaymentChannel, p.Status, p.SnapToken, p.SnapRedirectURL,
	).Scan(&p.ID, &p.UUID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *BillingRepository) GetLatestPaymentByInvoice(ctx context.Context, invoiceID int64) (*entity.Payment, error) {
	query := `
		SELECT id, uuid, payment_number, invoice_id, organization_id, amount,
		       COALESCE(payment_method, ''), COALESCE(payment_channel, ''), status,
		       COALESCE(snap_token, ''), COALESCE(snap_redirect_url, ''), paid_at, created_at, updated_at
		FROM payments
		WHERE invoice_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	p := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, invoiceID).
		Scan(&p.ID, &p.UUID, &p.PaymentNumber, &p.InvoiceID, &p.OrganizationID, &p.Amount,
			&p.PaymentMethod, &p.PaymentChannel, &p.Status,
			&p.SnapToken, &p.SnapRedirectURL, &p.PaidAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *BillingRepository) GetPaymentByNumber(ctx context.Context, paymentNum string) (*entity.Payment, error) {
	query := `
		SELECT id, uuid, payment_number, invoice_id, organization_id, amount,
		       COALESCE(payment_method, ''), COALESCE(payment_channel, ''), status,
		       COALESCE(snap_token, ''), COALESCE(snap_redirect_url, ''), paid_at, created_at, updated_at
		FROM payments
		WHERE payment_number = $1
	`
	p := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, paymentNum).
		Scan(&p.ID, &p.UUID, &p.PaymentNumber, &p.InvoiceID, &p.OrganizationID, &p.Amount,
			&p.PaymentMethod, &p.PaymentChannel, &p.Status,
			&p.SnapToken, &p.SnapRedirectURL, &p.PaidAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *BillingRepository) UpdatePaymentStatusTx(ctx context.Context, paymentNum string, paymentStatus string, method string, channel string, rawPayload string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update Payment Table
	var invoiceID int64
	var paidAt sql.NullTime
	if paymentStatus == "SETTLEMENT" || paymentStatus == "PAID" {
		paidAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	queryPayment := `
		UPDATE payments
		SET status = $1, payment_method = $2, payment_channel = $3, paid_at = $4, raw_gateway_response = $5::jsonb, updated_at = CURRENT_TIMESTAMP
		WHERE payment_number = $6
		RETURNING invoice_id
	`
	err = tx.QueryRowContext(ctx, queryPayment, paymentStatus, method, channel, paidAt, rawPayload, paymentNum).Scan(&invoiceID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// If Payment Settlement -> Update Invoice to PAID
	if paymentStatus == "SETTLEMENT" || paymentStatus == "PAID" {
		queryInv := `UPDATE invoices SET status = 'PAID', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
		_, err = tx.ExecContext(ctx, queryInv, invoiceID)
		if err != nil {
			return fmt.Errorf("failed to update invoice status to PAID: %w", err)
		}
	}

	return tx.Commit()
}

func (r *BillingRepository) SaveWebhookLog(ctx context.Context, orderID, status, payload string) error {
	query := `INSERT INTO billing_webhook_logs (order_id, transaction_status, raw_payload) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, orderID, status, payload)
	return err
}

func (r *BillingRepository) GetSuperadminStats(ctx context.Context) (*dto.SuperadminBillingStatsResponse, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'PAID' THEN total_amount ELSE 0 END), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN status = 'UNPAID' OR status = 'OVERDUE' THEN total_amount ELSE 0 END), 0) as pending_revenue,
			COUNT(CASE WHEN status = 'PAID' THEN 1 END) as paid_count,
			COUNT(CASE WHEN status = 'UNPAID' THEN 1 END) as unpaid_count,
			COUNT(CASE WHEN status = 'OVERDUE' THEN 1 END) as overdue_count
		FROM invoices
		WHERE deleted_at IS NULL
	`
	stats := &dto.SuperadminBillingStatsResponse{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&stats.TotalRevenue, &stats.PendingRevenue, &stats.PaidInvoicesCount, &stats.UnpaidInvoicesCount, &stats.OverdueInvoicesCount,
	)
	if err != nil {
		return nil, err
	}

	currentPeriod := time.Now().Format("2006-01")
	ticketQuery := `SELECT COALESCE(SUM(ticket_count), 0) FROM billing_usage_meters WHERE billing_period = $1`
	_ = r.db.QueryRowContext(ctx, ticketQuery, currentPeriod).Scan(&stats.CurrentMonthTickets)

	return stats, nil
}

func (r *BillingRepository) CheckAndSuspendOverdueOrganizations(ctx context.Context) error {
	markOverdueQuery := `
		UPDATE invoices
		SET status = 'OVERDUE', updated_at = CURRENT_TIMESTAMP
		WHERE status = 'UNPAID' AND due_date < CURRENT_TIMESTAMP AND deleted_at IS NULL
	`
	_, _ = r.db.ExecContext(ctx, markOverdueQuery)

	suspendOrgQuery := `
		UPDATE organizations
		SET status = 'SUSPENDED', updated_at = CURRENT_TIMESTAMP
		WHERE id IN (
			SELECT organization_id
			FROM invoices
			WHERE status = 'OVERDUE' AND due_date < (CURRENT_TIMESTAMP - INTERVAL '7 days') AND deleted_at IS NULL
		) AND status = 'ACTIVE'
	`
	_, err := r.db.ExecContext(ctx, suspendOrgQuery)
	return err
}
