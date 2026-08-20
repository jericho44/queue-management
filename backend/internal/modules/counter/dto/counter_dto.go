package dto

type CreateCounterRequest struct {
	BranchID      int64   `json:"branch_id" validate:"required"`
	CounterNumber string  `json:"counter_number" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	ServiceIDs    []int64 `json:"service_ids"`
}

type CounterAssignStaffRequest struct {
	StaffID int64 `json:"staff_id"`
}
