package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int64          `json:"id"`
	UUID           string         `json:"uuid"`
	OrganizationID sql.NullInt64  `json:"organization_id,omitempty"`
	OrgUUID        string         `json:"org_uuid,omitempty"`
	OrgName        string         `json:"org_name,omitempty"`
	Email          string         `json:"email"`
	PasswordHash   string         `json:"-"`
	FullName       string         `json:"full_name"`
	Role           string         `json:"role"`
	Phone          sql.NullString `json:"phone,omitempty"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at,omitempty"`
}
