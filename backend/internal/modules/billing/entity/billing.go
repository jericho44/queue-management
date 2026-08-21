package entity

import (
	"database/sql"
	"time"
)

type BillingPlan struct {
	ID             int64     `json:"id"`
	PlanCode       string    `json:"plan_code"`
	Name           string    `json:"name"`
	RatePerTicket  float64   `json:"rate_per_ticket"`
	BaseMonthlyFee float64   `json:"base_monthly_fee"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UsageMeter struct {
	ID             int64        `json:"id"`
	OrganizationID int64        `json:"organization_id"`
	BillingPeriod  string       `json:"billing_period"`
	TicketCount    int          `json:"ticket_count"`
	LastTicketAt   sql.NullTime `json:"last_ticket_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type InvoiceItem struct {
	ID          int64     `json:"id"`
	InvoiceID   int64     `json:"invoice_id"`
	Description string    `json:"description"`
	ItemType    string    `json:"item_type"` // TICKET_USAGE, BASE_FEE, ADDON
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type Invoice struct {
	ID             int64         `json:"id"`
	UUID           string        `json:"uuid"`
	InvoiceNumber  string        `json:"invoice_number"`
	OrganizationID int64         `json:"organization_id"`
	OrgName        string        `json:"org_name,omitempty"`
	BillingPeriod  string        `json:"billing_period"`
	Subtotal       float64       `json:"subtotal"`
	TaxAmount      float64       `json:"tax_amount"`
	TotalAmount    float64       `json:"total_amount"`
	Status         string        `json:"status"` // UNPAID, PAID, OVERDUE, CANCELLED
	IssuedAt       time.Time     `json:"issued_at"`
	DueDate        time.Time     `json:"due_date"`
	Items          []InvoiceItem `json:"items,omitempty"`
	LatestPayment  *Payment      `json:"latest_payment,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type Payment struct {
	ID                 int64        `json:"id"`
	UUID               string       `json:"uuid"`
	PaymentNumber      string       `json:"payment_number"`
	InvoiceID          int64        `json:"invoice_id"`
	OrganizationID     int64        `json:"organization_id"`
	Amount             float64      `json:"amount"`
	PaymentMethod      string       `json:"payment_method,omitempty"`
	PaymentChannel     string       `json:"payment_channel,omitempty"`
	Status             string       `json:"status"` // PENDING, SETTLEMENT, DENIED, EXPIRED, CANCELLED
	SnapToken          string       `json:"snap_token,omitempty"`
	SnapRedirectURL   string       `json:"snap_redirect_url,omitempty"`
	PaidAt             sql.NullTime `json:"paid_at,omitempty"`
	RawGatewayResponse string       `json:"raw_gateway_response,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

type WebhookLog struct {
	ID                int64     `json:"id"`
	OrderID           string    `json:"order_id"`
	TransactionStatus string    `json:"transaction_status"`
	RawPayload        string    `json:"raw_payload"`
	CreatedAt         time.Time `json:"created_at"`
}
