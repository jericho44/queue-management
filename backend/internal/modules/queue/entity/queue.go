package entity

import (
	"database/sql"
	"time"
)

type QueueTicket struct {
	ID                   int64        `json:"id"`
	UUID                 string       `json:"uuid"`
	OrganizationID       int64        `json:"organization_id"`
	BranchID             int64        `json:"branch_id"`
	BranchName           string       `json:"branch_name,omitempty"`
	ServiceID            int64        `json:"service_id"`
	ServiceName          string       `json:"service_name,omitempty"`
	ServicePrefix        string       `json:"service_prefix,omitempty"`
	CustomerID           sql.NullInt64 `json:"customer_id,omitempty"`
	CustomerName         string       `json:"customer_name,omitempty"`
	CounterID            sql.NullInt64 `json:"counter_id,omitempty"`
	CounterNumber        string       `json:"counter_number,omitempty"`
	StaffID              sql.NullInt64 `json:"staff_id,omitempty"`
	StaffName            string       `json:"staff_name,omitempty"`
	TicketNumber         string       `json:"ticket_number"`
	SequenceNumber       int          `json:"sequence_number"`
	QueueDate            time.Time    `json:"queue_date"`
	Priority             string       `json:"priority"` // NORMAL, PRIORITY, EMERGENCY
	Status               string       `json:"status"`   // WAITING, CALLED, SERVING, COMPLETED, SKIPPED, CANCELLED, NO_SHOW, TRANSFERRED
	PublicToken          string       `json:"public_token"`
	CalledAt             sql.NullTime `json:"called_at,omitempty"`
	ServingStartedAt     sql.NullTime `json:"serving_started_at,omitempty"`
	CompletedAt          sql.NullTime `json:"completed_at,omitempty"`
	CancelledAt          sql.NullTime `json:"cancelled_at,omitempty"`
	EstimatedWaitSeconds int          `json:"estimated_wait_seconds"`
	PeopleAhead          int          `json:"people_ahead,omitempty"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	DeletedAt            sql.NullTime `json:"deleted_at,omitempty"`
}

type QueueTicketEvent struct {
	ID         int64        `json:"id"`
	TicketID   int64        `json:"ticket_id"`
	FromStatus string       `json:"from_status"`
	ToStatus   string       `json:"to_status"`
	ActorID    sql.NullInt64 `json:"actor_id,omitempty"`
	ActorName  string       `json:"actor_name,omitempty"`
	Metadata   string       `json:"metadata"`
	CreatedAt  time.Time    `json:"created_at"`
}
