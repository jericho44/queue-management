package entity

import (
	"database/sql"
	"time"
)

type AuditLog struct {
	ID             int64        `json:"id"`
	UUID           string       `json:"uuid"`
	OrganizationID sql.NullInt64 `json:"organization_id,omitempty"`
	BranchID       sql.NullInt64 `json:"branch_id,omitempty"`
	UserID         sql.NullInt64 `json:"user_id,omitempty"`
	UserName       string       `json:"user_name,omitempty"`
	Action         string       `json:"action"`
	EntityType     string       `json:"entity_type"`
	EntityID       sql.NullInt64 `json:"entity_id,omitempty"`
	OldValues      string       `json:"old_values,omitempty"`
	NewValues      string       `json:"new_values,omitempty"`
	IPAddress      string       `json:"ip_address,omitempty"`
	UserAgent      string       `json:"user_agent,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}
