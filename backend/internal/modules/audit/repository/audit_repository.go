package repository

import (
	"context"
	"database/sql"

	"queue-management-tenant/backend/internal/modules/audit/entity"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) LogAction(ctx context.Context, log *entity.AuditLog) error {
	query := `
		INSERT INTO audit_logs (organization_id, branch_id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	var orgID, branchID, userID, entityID interface{} = nil, nil, nil, nil
	if log.OrganizationID.Valid {
		orgID = log.OrganizationID.Int64
	}
	if log.BranchID.Valid {
		branchID = log.BranchID.Int64
	}
	if log.UserID.Valid {
		userID = log.UserID.Int64
	}
	if log.EntityID.Valid {
		entityID = log.EntityID.Int64
	}

	_, err := r.db.ExecContext(ctx, query,
		orgID, branchID, userID, log.Action, log.EntityType, entityID, log.OldValues, log.NewValues, log.IPAddress, log.UserAgent,
	)
	return err
}

func (r *AuditRepository) List(ctx context.Context, orgID int64, page, perPage int) ([]entity.AuditLog, int64, error) {
	countQuery := `SELECT COUNT(*) FROM audit_logs WHERE organization_id = $1`
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, orgID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	query := `
		SELECT a.id, a.uuid, a.organization_id, a.branch_id, a.user_id, COALESCE(u.full_name, ''), a.action, a.entity_type, a.entity_id, COALESCE(a.old_values::text, ''), COALESCE(a.new_values::text, ''), COALESCE(a.ip_address, ''), COALESCE(a.user_agent, ''), a.created_at
		FROM audit_logs a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.organization_id = $1
		ORDER BY a.id DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, orgID, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []entity.AuditLog
	for rows.Next() {
		var l entity.AuditLog
		if err := rows.Scan(
			&l.ID, &l.UUID, &l.OrganizationID, &l.BranchID, &l.UserID, &l.UserName,
			&l.Action, &l.EntityType, &l.EntityID, &l.OldValues, &l.NewValues, &l.IPAddress, &l.UserAgent, &l.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, nil
}
