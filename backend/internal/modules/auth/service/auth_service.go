package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"queue-management-tenant/backend/internal/modules/auth/dto"
	authEntity "queue-management-tenant/backend/internal/modules/auth/entity"
	authRepo "queue-management-tenant/backend/internal/modules/auth/repository"
	orgEntity "queue-management-tenant/backend/internal/modules/organization/entity"
	orgRepo "queue-management-tenant/backend/internal/modules/organization/repository"
	"queue-management-tenant/backend/pkg/jwt"
)

type AuthService struct {
	db       *sql.DB
	orgRepo  *orgRepo.OrganizationRepository
	userRepo *authRepo.UserRepository
	jwtSvc   *jwt.JWTService
}

func NewAuthService(db *sql.DB, orgRepo *orgRepo.OrganizationRepository, userRepo *authRepo.UserRepository, jwtSvc *jwt.JWTService) *AuthService {
	return &AuthService{
		db:       db,
		orgRepo:  orgRepo,
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

func (s *AuthService) RegisterOrganization(ctx context.Context, req dto.RegisterOrganizationRequest) (*dto.LoginResponse, error) {
	slug := strings.ToLower(strings.ReplaceAll(req.OrgCode, " ", "-"))

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Create Organization
	org := &orgEntity.Organization{
		Name:     req.OrgName,
		Code:     strings.ToUpper(req.OrgCode),
		Slug:     slug,
		Status:   "ACTIVE",
		Settings: "{}",
	}
	if err := s.orgRepo.CreateTx(ctx, tx, org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// 2. Create Subscription
	sub := &orgEntity.Subscription{
		OrganizationID:    org.ID,
		Plan:              "FREE",
		MaxBranches:       1,
		MaxCounters:       2,
		MaxStaff:          5,
		MaxMonthlyTickets: 1000,
		Status:            "ACTIVE",
	}
	if err := s.orgRepo.CreateSubscriptionTx(ctx, tx, sub); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// 3. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Create Owner User
	user := &authEntity.User{
		OrganizationID: sql.NullInt64{Int64: org.ID, Valid: true},
		Email:          req.AdminEmail,
		PasswordHash:   string(hashed),
		FullName:       req.AdminName,
		Role:           "OWNER",
		Status:         "ACTIVE",
	}
	if err := s.userRepo.CreateTx(ctx, tx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Generate JWT
	accessTok, err := s.jwtSvc.GenerateAccessToken(user.ID, user.UUID, org.ID, org.UUID, user.Email, user.FullName, user.Role)
	if err != nil {
		return nil, err
	}
	refreshTok, err := s.jwtSvc.GenerateRefreshToken(user.ID, user.UUID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
		User:         user,
		Organization: org,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	var orgID int64 = 0
	if user.OrganizationID.Valid {
		orgID = user.OrganizationID.Int64
	}

	accessTok, err := s.jwtSvc.GenerateAccessToken(user.ID, user.UUID, orgID, user.OrgUUID, user.Email, user.FullName, user.Role)
	if err != nil {
		return nil, err
	}

	refreshTok, err := s.jwtSvc.GenerateRefreshToken(user.ID, user.UUID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
		User:         user,
	}, nil
}

func (s *AuthService) CreateUser(ctx context.Context, orgID int64, req dto.CreateUserRequest) (*authEntity.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &authEntity.User{
		OrganizationID: sql.NullInt64{Int64: orgID, Valid: true},
		Email:          req.Email,
		PasswordHash:   string(hashed),
		FullName:       req.FullName,
		Role:           req.Role,
		Phone:          sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Status:         "ACTIVE",
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
