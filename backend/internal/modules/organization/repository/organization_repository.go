package repository

import (
	"context"
	"database/sql"

	"queue-management-tenant/backend/internal/modules/organization/entity"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) CreateTx(ctx context.Context, tx *sql.Tx, org *entity.Organization) error {
	query := `
		INSERT INTO organizations (name, code, slug, status, settings)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, created_at, updated_at
	`
	return tx.QueryRowContext(ctx, query, org.Name, org.Code, org.Slug, org.Status, org.Settings).
		Scan(&org.ID, &org.UUID, &org.CreatedAt, &org.UpdatedAt)
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id int64) (*entity.Organization, error) {
	query := `
		SELECT id, uuid, name, code, slug, status, settings, created_at, updated_at
		FROM organizations WHERE id = $1 AND deleted_at IS NULL
	`
	org := &entity.Organization{}
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&org.ID, &org.UUID, &org.Name, &org.Code, &org.Slug, &org.Status, &org.Settings, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (r *OrganizationRepository) CreateSubscriptionTx(ctx context.Context, tx *sql.Tx, sub *entity.Subscription) error {
	query := `
		INSERT INTO subscriptions (organization_id, plan, max_branches, max_counters, max_staff, max_monthly_tickets, status, starts_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, uuid, created_at, updated_at
	`
	return tx.QueryRowContext(ctx, query,
		sub.OrganizationID, sub.Plan, sub.MaxBranches, sub.MaxCounters, sub.MaxStaff, sub.MaxMonthlyTickets, sub.Status, sub.StartsAt, sub.ExpiresAt,
	).Scan(&sub.ID, &sub.UUID, &sub.CreatedAt, &sub.UpdatedAt)
}

func (r *OrganizationRepository) GetActiveSubscription(ctx context.Context, orgID int64) (*entity.Subscription, error) {
	query := `
		SELECT id, uuid, organization_id, plan, max_branches, max_counters, max_staff, max_monthly_tickets, status, starts_at, expires_at, created_at, updated_at
		FROM subscriptions
		WHERE organization_id = $1 AND status = 'ACTIVE' AND deleted_at IS NULL
		ORDER BY id DESC LIMIT 1
	`
	sub := &entity.Subscription{}
	err := r.db.QueryRowContext(ctx, query, orgID).
		Scan(&sub.ID, &sub.UUID, &sub.OrganizationID, &sub.Plan, &sub.MaxBranches, &sub.MaxCounters, &sub.MaxStaff, &sub.MaxMonthlyTickets, &sub.Status, &sub.StartsAt, &sub.ExpiresAt, &sub.CreatedAt, &sub.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return sub, nil
}
