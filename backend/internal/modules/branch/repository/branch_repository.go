package repository

import (
	"context"
	"database/sql"

	"queue-management-tenant/backend/internal/modules/branch/entity"
)

type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(ctx context.Context, b *entity.Branch) error {
	query := `
		INSERT INTO branches (organization_id, name, code, address, phone, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uuid, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		b.OrganizationID, b.Name, b.Code, b.Address, b.Phone, b.Status,
	).Scan(&b.ID, &b.UUID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *BranchRepository) GetByID(ctx context.Context, orgID, id int64) (*entity.Branch, error) {
	query := `
		SELECT id, uuid, organization_id, name, code, address, phone, status, created_at, updated_at
		FROM branches WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
	b := &entity.Branch{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).
		Scan(&b.ID, &b.UUID, &b.OrganizationID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Status, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BranchRepository) ListByOrg(ctx context.Context, orgID int64) ([]entity.Branch, error) {
	query := `
		SELECT id, uuid, organization_id, name, code, address, phone, status, created_at, updated_at
		FROM branches WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []entity.Branch
	for rows.Next() {
		var b entity.Branch
		if err := rows.Scan(&b.ID, &b.UUID, &b.OrganizationID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Status, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (r *BranchRepository) CountByOrg(ctx context.Context, orgID int64) (int, error) {
	query := `SELECT COUNT(*) FROM branches WHERE organization_id = $1 AND deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(&count)
	return count, err
}
