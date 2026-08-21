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
	Address        string       `json:"address"`
	Phone          string       `json:"phone"`
	Status         string       `json:"status"`
	KioskEnabled  bool         `json:"kiosk_enabled"`
	KioskMode     string       `json:"kiosk_mode"`
	PaperSize     string       `json:"paper_size"`
	ReceiptHeader string       `json:"receipt_header"`
	ReceiptFooter string       `json:"receipt_footer"`
	AutoPrint     bool         `json:"auto_print"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty"`
}
