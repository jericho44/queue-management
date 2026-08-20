package dto

type CreateBranchRequest struct {
	Name    string `json:"name" validate:"required"`
	Code    string `json:"code" validate:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}
