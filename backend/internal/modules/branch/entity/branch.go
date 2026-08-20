package entity

import (
	"database/sql"
	"time"
)

type Branch struct {
	ID             int64        `json:"id"`
	UUID           string       `json:"uuid"`
	OrganizationID int64        `json:"organization_id"`
	Name           string       `json:"name"`
	Code           string       `json:"code"`
	Address        sql.NullString `json:"address,omitempty"`
	Phone          sql.NullString `json:"phone,omitempty"`
	Status         string       `json:"status"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty"`
}
