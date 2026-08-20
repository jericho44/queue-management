package entity

import (
	"database/sql"
	"time"
)

type Service struct {
	ID               int64        `json:"id"`
	UUID             string       `json:"uuid"`
	OrganizationID   int64        `json:"organization_id"`
	BranchID         int64        `json:"branch_id"`
	BranchName       string       `json:"branch_name,omitempty"`
	Name             string       `json:"name"`
	Code             string       `json:"code"`
	Prefix           string       `json:"prefix"`
	AvgServiceTimeSec int         `json:"avg_service_time_sec"`
	PriorityWeight   int          `json:"priority_weight"`
	Status           string       `json:"status"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	DeletedAt        sql.NullTime `json:"deleted_at,omitempty"`
}
