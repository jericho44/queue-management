package dto

type CreateServiceRequest struct {
	BranchID          int64  `json:"branch_id" validate:"required"`
	Name              string `json:"name" validate:"required"`
	Code              string `json:"code" validate:"required"`
	Prefix            string `json:"prefix" validate:"required"`
	AvgServiceTimeSec int    `json:"avg_service_time_sec"`
	PriorityWeight    int    `json:"priority_weight"`
}
