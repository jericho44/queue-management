package repository

import (
	"context"
	"database/sql"

	counterEntity "queue-management-tenant/backend/internal/modules/counter/entity"
	serviceEntity "queue-management-tenant/backend/internal/modules/service/entity"
)

type CounterRepository struct {
	db *sql.DB
}

func NewCounterRepository(db *sql.DB) *CounterRepository {
	return &CounterRepository{db: db}
}

func (r *CounterRepository) Create(ctx context.Context, c *counterEntity.Counter, serviceIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO counters (organization_id, branch_id, counter_number, name, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, uuid, created_at, updated_at
	`
	err = tx.QueryRowContext(ctx, query,
		c.OrganizationID, c.BranchID, c.CounterNumber, c.Name, c.Status,
	).Scan(&c.ID, &c.UUID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return err
	}

	for _, sID := range serviceIDs {
		_, err := tx.ExecContext(ctx, `INSERT INTO counter_services (counter_id, service_id) VALUES ($1, $2)`, c.ID, sID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *CounterRepository) GetByID(ctx context.Context, orgID, id int64) (*counterEntity.Counter, error) {
	query := `
		SELECT c.id, c.uuid, c.organization_id, c.branch_id, b.name, c.counter_number, c.name, c.status, c.current_staff_id, COALESCE(u.full_name, ''), c.created_at, c.updated_at
		FROM counters c
		JOIN branches b ON c.branch_id = b.id
		LEFT JOIN users u ON c.current_staff_id = u.id
		WHERE c.id = $1 AND c.organization_id = $2 AND c.deleted_at IS NULL
	`
	c := &counterEntity.Counter{}
	err := r.db.QueryRowContext(ctx, query, id, orgID).
		Scan(&c.ID, &c.UUID, &c.OrganizationID, &c.BranchID, &c.BranchName, &c.CounterNumber, &c.Name, &c.Status, &c.CurrentStaffID, &c.StaffName, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	c.Services, _ = r.getCounterServices(ctx, c.ID)
	return c, nil
}

func (r *CounterRepository) getCounterServices(ctx context.Context, counterID int64) ([]serviceEntity.Service, error) {
	query := `
		SELECT s.id, s.uuid, s.organization_id, s.branch_id, s.name, s.code, s.prefix, s.avg_service_time_sec, s.priority_weight, s.status, s.created_at, s.updated_at
		FROM services s
		JOIN counter_services cs ON s.id = cs.service_id
		WHERE cs.counter_id = $1 AND s.deleted_at IS NULL
	`
	rows, err := r.db.QueryContext(ctx, query, counterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []serviceEntity.Service
	for rows.Next() {
		var s serviceEntity.Service
		if err := rows.Scan(&s.ID, &s.UUID, &s.OrganizationID, &s.BranchID, &s.Name, &s.Code, &s.Prefix, &s.AvgServiceTimeSec, &s.PriorityWeight, &s.Status, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}

func (r *CounterRepository) ListByBranch(ctx context.Context, orgID, branchID int64) ([]counterEntity.Counter, error) {
	query := `
		SELECT c.id, c.uuid, c.organization_id, c.branch_id, b.name, c.counter_number, c.name, c.status, c.current_staff_id, COALESCE(u.full_name, ''), c.created_at, c.updated_at
		FROM counters c
		JOIN branches b ON c.branch_id = b.id
		LEFT JOIN users u ON c.current_staff_id = u.id
		WHERE c.organization_id = $1 AND c.branch_id = $2 AND c.deleted_at IS NULL
		ORDER BY c.counter_number ASC
	`
	rows, err := r.db.QueryContext(ctx, query, orgID, branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counters []counterEntity.Counter
	for rows.Next() {
		var c counterEntity.Counter
		if err := rows.Scan(&c.ID, &c.UUID, &c.OrganizationID, &c.BranchID, &c.BranchName, &c.CounterNumber, &c.Name, &c.Status, &c.CurrentStaffID, &c.StaffName, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		c.Services, _ = r.getCounterServices(ctx, c.ID)
		counters = append(counters, c)
	}
	return counters, nil
}

func (r *CounterRepository) UpdateStatusAndStaff(ctx context.Context, counterID int64, status string, staffID sql.NullInt64) error {
	query := `UPDATE counters SET status = $1, current_staff_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	var sID interface{} = nil
	if staffID.Valid {
		sID = staffID.Int64
	}
	_, err := r.db.ExecContext(ctx, query, status, sID, counterID)
	return err
}

func (r *CounterRepository) CountByOrg(ctx context.Context, orgID int64) (int, error) {
	query := `SELECT COUNT(*) FROM counters WHERE organization_id = $1 AND deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(&count)
	return count, err
}
