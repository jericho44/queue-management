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
		INSERT INTO branches (organization_id, name, code, address, phone, status, kiosk_enabled, kiosk_mode, paper_size, receipt_header, receipt_footer, auto_print)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, uuid, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		b.OrganizationID, b.Name, b.Code, b.Address, b.Phone, b.Status,
		b.KioskEnabled, b.KioskMode, b.PaperSize, b.ReceiptHeader, b.ReceiptFooter, b.AutoPrint,
	).Scan(&b.ID, &b.UUID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *BranchRepository) GetByID(ctx context.Context, orgID, id int64) (*entity.Branch, error) {
	query := `
		SELECT id, uuid, organization_id, name, code, address, phone, status,
		       kiosk_enabled, kiosk_mode, paper_size, receipt_header, receipt_footer, auto_print,
		       created_at, updated_at
		FROM branches WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`
	b := &entity.Branch{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).
		Scan(&b.ID, &b.UUID, &b.OrganizationID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Status,
			&b.KioskEnabled, &b.KioskMode, &b.PaperSize, &b.ReceiptHeader, &b.ReceiptFooter, &b.AutoPrint,
			&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BranchRepository) GetByIDPublic(ctx context.Context, id int64) (*entity.Branch, error) {
	query := `
		SELECT id, uuid, organization_id, name, code, address, phone, status,
		       kiosk_enabled, kiosk_mode, paper_size, receipt_header, receipt_footer, auto_print,
		       created_at, updated_at
		FROM branches WHERE id = $1 AND deleted_at IS NULL
	`
	b := &entity.Branch{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&b.ID, &b.UUID, &b.OrganizationID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Status,
			&b.KioskEnabled, &b.KioskMode, &b.PaperSize, &b.ReceiptHeader, &b.ReceiptFooter, &b.AutoPrint,
			&b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r *BranchRepository) ListByOrg(ctx context.Context, orgID int64) ([]entity.Branch, error) {
	query := `
		SELECT id, uuid, organization_id, name, code, address, phone, status,
		       kiosk_enabled, kiosk_mode, paper_size, receipt_header, receipt_footer, auto_print,
		       created_at, updated_at
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
		if err := rows.Scan(&b.ID, &b.UUID, &b.OrganizationID, &b.Name, &b.Code, &b.Address, &b.Phone, &b.Status,
			&b.KioskEnabled, &b.KioskMode, &b.PaperSize, &b.ReceiptHeader, &b.ReceiptFooter, &b.AutoPrint,
			&b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (r *BranchRepository) UpdateKioskSettings(ctx context.Context, b *entity.Branch) error {
	query := `
		UPDATE branches
		SET kiosk_enabled = $1, kiosk_mode = $2, paper_size = $3, receipt_header = $4, receipt_footer = $5, auto_print = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND organization_id = $8 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query,
		b.KioskEnabled, b.KioskMode, b.PaperSize, b.ReceiptHeader, b.ReceiptFooter, b.AutoPrint,
		b.ID, b.OrganizationID,
	)
	return err
}

func (r *BranchRepository) CountByOrg(ctx context.Context, orgID int64) (int, error) {
	query := `SELECT COUNT(*) FROM branches WHERE organization_id = $1 AND deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(&count)
	return count, err
}

