package repository

import (
	"context"
	"database/sql"

	"queue-management-tenant/backend/internal/modules/auth/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *entity.User) error {
	query := `
		INSERT INTO users (organization_id, email, password_hash, full_name, role, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, created_at, updated_at
	`
	var orgID interface{} = nil
	if u.OrganizationID.Valid {
		orgID = u.OrganizationID.Int64
	}

	return r.db.QueryRowContext(ctx, query,
		orgID, u.Email, u.PasswordHash, u.FullName, u.Role, u.Phone, u.Status,
	).Scan(&u.ID, &u.UUID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) CreateTx(ctx context.Context, tx *sql.Tx, u *entity.User) error {
	query := `
		INSERT INTO users (organization_id, email, password_hash, full_name, role, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, created_at, updated_at
	`
	var orgID interface{} = nil
	if u.OrganizationID.Valid {
		orgID = u.OrganizationID.Int64
	}

	return tx.QueryRowContext(ctx, query,
		orgID, u.Email, u.PasswordHash, u.FullName, u.Role, u.Phone, u.Status,
	).Scan(&u.ID, &u.UUID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT u.id, u.uuid, u.organization_id, COALESCE(o.uuid::text, ''), COALESCE(o.name, ''), u.email, u.password_hash, u.full_name, u.role, u.phone, u.status, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN organizations o ON u.organization_id = o.id
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`
	u := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.UUID, &u.OrganizationID, &u.OrgUUID, &u.OrgName, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.Phone, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT u.id, u.uuid, u.organization_id, COALESCE(o.uuid::text, ''), COALESCE(o.name, ''), u.email, u.password_hash, u.full_name, u.role, u.phone, u.status, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN organizations o ON u.organization_id = o.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`
	u := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&u.ID, &u.UUID, &u.OrganizationID, &u.OrgUUID, &u.OrgName, &u.Email, &u.PasswordHash, &u.FullName, &u.Role, &u.Phone, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
