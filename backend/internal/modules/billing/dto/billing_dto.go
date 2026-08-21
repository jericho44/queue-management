package dto

type MidtransWebhookPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type CreateSnapTokenRequest struct {
	InvoiceID int64 `json:"invoice_id" validate:"required"`
}

type CreateSnapTokenResponse struct {
	PaymentNumber   string `json:"payment_number"`
	SnapToken       string `json:"snap_token"`
	SnapRedirectURL string `json:"snap_redirect_url"`
}

type SuperadminBillingStatsResponse struct {
	TotalRevenue        float64 `json:"total_revenue"`
	PendingRevenue      float64 `json:"pending_revenue"`
	PaidInvoicesCount   int     `json:"paid_invoices_count"`
	UnpaidInvoicesCount int     `json:"unpaid_invoices_count"`
	OverdueInvoicesCount int    `json:"overdue_invoices_count"`
	CurrentMonthTickets int     `json:"current_month_tickets"`
}
