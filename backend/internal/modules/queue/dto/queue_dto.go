package dto

type IssueTicketRequest struct {
	BranchID  int64  `json:"branch_id" validate:"required"`
	ServiceID int64  `json:"service_id" validate:"required"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Priority  string `json:"priority"` // NORMAL, PRIORITY, EMERGENCY
}

type TransferTicketRequest struct {
	TargetServiceID int64  `json:"target_service_id" validate:"required"`
	Reason          string `json:"reason"`
}
