package dto

type CreateBranchRequest struct {
	Name    string `json:"name" validate:"required"`
	Code    string `json:"code" validate:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UpdateKioskSettingsRequest struct {
	KioskEnabled  *bool   `json:"kiosk_enabled"`
	KioskMode     *string `json:"kiosk_mode"`
	PaperSize     *string `json:"paper_size"`
	ReceiptHeader *string `json:"receipt_header"`
	ReceiptFooter *string `json:"receipt_footer"`
	AutoPrint     *bool   `json:"auto_print"`
}

