package dto

type RegisterOrganizationRequest struct {
	OrgName    string `json:"org_name" validate:"required"`
	OrgCode    string `json:"org_code" validate:"required"`
	AdminEmail string `json:"admin_email" validate:"required,email"`
	AdminName  string `json:"admin_name" validate:"required"`
	Password   string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         interface{} `json:"user"`
	Organization interface{} `json:"organization,omitempty"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required"`
	Phone    string `json:"phone"`
}
