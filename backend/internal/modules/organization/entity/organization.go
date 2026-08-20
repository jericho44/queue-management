package entity

import (
	"database/sql"
	"time"
)

type Organization struct {
	ID        int64          `json:"id"`
	UUID      string         `json:"uuid"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Slug      string         `json:"slug"`
	Status    string         `json:"status"`
	Settings  string         `json:"settings"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt sql.NullTime   `json:"deleted_at,omitempty"`
	CreatedBy sql.NullInt64  `json:"created_by,omitempty"`
	UpdatedBy sql.NullInt64  `json:"updated_by,omitempty"`
	DeletedBy sql.NullInt64  `json:"deleted_by,omitempty"`
}

type Subscription struct {
	ID                int64        `json:"id"`
	UUID              string       `json:"uuid"`
	OrganizationID    int64        `json:"organization_id"`
	Plan              string       `json:"plan"`
	MaxBranches       int          `json:"max_branches"`
	MaxCounters       int          `json:"max_counters"`
	MaxStaff          int          `json:"max_staff"`
	MaxMonthlyTickets int          `json:"max_monthly_tickets"`
	Status            string       `json:"status"`
	StartsAt          time.Time    `json:"starts_at"`
	ExpiresAt         sql.NullTime `json:"expires_at,omitempty"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty"`
}
