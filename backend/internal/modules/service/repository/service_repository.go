package repository

import (
	"context"
	"database/sql"

	"queue-management-tenant/backend/internal/modules/service/entity"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(ctx context.Context, s *entity.Service) error {
	query := `
		INSERT INTO services (organization_id, branch_id, name, code, prefix, avg_service_time_sec, priority_weight, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, uuid, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		s.OrganizationID, s.BranchID, s.Name, s.Code, s.Prefix, s.AvgServiceTimeSec, s.PriorityWeight, s.Status,
	).Scan(&s.ID, &s.UUID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *ServiceRepository) GetByID(ctx context.Context, orgID, id int64) (*entity.Service, error) {
	query := `
		SELECT s.id, s.uuid, s.organization_id, s.branch_id, b.name, s.name, s.code, s.prefix, s.avg_service_time_sec, s.priority_weight, s.status, s.created_at, s.updated_at
		FROM services s
		JOIN branches b ON s.branch_id = b.id
		WHERE s.id = $1 AND s.organization_id = $2 AND s.deleted_at IS NULL
	`
	s := &entity.Service{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).
		Scan(&s.ID, &s.UUID, &s.OrganizationID, &s.BranchID, &s.BranchName, &s.Name, &s.Code, &s.Prefix, &s.AvgServiceTimeSec, &s.PriorityWeight, &s.Status, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *ServiceRepository) ListByBranch(ctx context.Context, orgID, branchID int64) ([]entity.Service, error) {
	query := `
		SELECT s.id, s.uuid, s.organization_id, s.branch_id, b.name, s.name, s.code, s.prefix, s.avg_service_time_sec, s.priority_weight, s.status, s.created_at, s.updated_at
		FROM services s
		JOIN branches b ON s.branch_id = b.id
		WHERE s.organization_id = $1 AND s.branch_id = $2 AND s.deleted_at IS NULL
		ORDER BY s.prefix ASC, s.name ASC
	`
	rows, err := r.db.QueryContext(ctx, query, orgID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var s entity.Service
		if err := rows.Scan(&s.ID, &s.UUID, &s.OrganizationID, &s.BranchID, &s.BranchName, &s.Name, &s.Code, &s.Prefix, &s.AvgServiceTimeSec, &s.PriorityWeight, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
