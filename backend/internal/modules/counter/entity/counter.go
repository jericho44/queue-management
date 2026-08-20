package entity

import (
	"database/sql"
	"time"

	serviceEntity "queue-management-tenant/backend/internal/modules/service/entity"
)

type Counter struct {
	ID             int64                  `json:"id"`
	UUID           string                 `json:"uuid"`
	OrganizationID int64                  `json:"organization_id"`
	BranchID       int64                  `json:"branch_id"`
	BranchName     string                 `json:"branch_name,omitempty"`
	CounterNumber  string                 `json:"counter_number"`
	Name           string                 `json:"name"`
	Status         string                 `json:"status"` // CLOSED, OPEN, BUSY, PAUSED
	CurrentStaffID sql.NullInt64          `json:"current_staff_id,omitempty"`
	StaffName      string                 `json:"staff_name,omitempty"`
	Services       []serviceEntity.Service `json:"services,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	DeletedAt      sql.NullTime           `json:"deleted_at,omitempty"`
}
