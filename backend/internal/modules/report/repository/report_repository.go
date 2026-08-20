package repository

import (
	"context"
	"database/sql"
	"time"

	"queue-management-tenant/backend/internal/modules/report/dto"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetDashboardStats(ctx context.Context, orgID, branchID int64) (*dto.DashboardStatsResponse, error) {
	today := time.Now().UTC().Format("2006-01-02")

	stats := &dto.DashboardStatsResponse{}

	queryCounts := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'COMPLETED'),
			COUNT(*) FILTER (WHERE status = 'WAITING'),
			COUNT(*) FILTER (WHERE status IN ('CALLED', 'SERVING')),
			COUNT(*) FILTER (WHERE status = 'NO_SHOW'),
			COUNT(*) FILTER (WHERE status = 'CANCELLED')
		FROM queue_tickets
		WHERE organization_id = $1 AND ($2::bigint = 0 OR branch_id = $2) AND queue_date = $3 AND deleted_at IS NULL
	`
	_ = r.db.QueryRowContext(ctx, queryCounts, orgID, branchID, today).Scan(
		&stats.TotalTicketsToday, &stats.CompletedToday, &stats.WaitingCount, &stats.ServingCount, &stats.NoShowCount, &stats.CancelledCount,
	)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM counters
		WHERE organization_id = $1 AND ($2::bigint = 0 OR branch_id = $2) AND status IN ('OPEN', 'BUSY') AND deleted_at IS NULL
	`, orgID, branchID).Scan(&stats.ActiveCounters)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (called_at - created_at))), 0)::bigint
		FROM queue_tickets
		WHERE organization_id = $1 AND ($2::bigint = 0 OR branch_id = $2) AND queue_date = $3 AND called_at IS NOT NULL AND deleted_at IS NULL
	`, orgID, branchID, today).Scan(&stats.AvgWaitTimeSec)

	_ = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (completed_at - serving_started_at))), 0)::bigint
		FROM queue_tickets
		WHERE organization_id = $1 AND ($2::bigint = 0 OR branch_id = $2) AND queue_date = $3 AND completed_at IS NOT NULL AND serving_started_at IS NOT NULL AND deleted_at IS NULL
	`, orgID, branchID, today).Scan(&stats.AvgServiceTimeSec)

	hourlyQuery := `
		SELECT EXTRACT(HOUR FROM created_at)::int AS hr, COUNT(*)
		FROM queue_tickets
		WHERE organization_id = $1 AND ($2::bigint = 0 OR branch_id = $2) AND queue_date = $3 AND deleted_at IS NULL
		GROUP BY hr ORDER BY hr ASC
	`
	rowsH, err := r.db.QueryContext(ctx, hourlyQuery, orgID, branchID, today)
	if err == nil {
		defer rowsH.Close()
		for rowsH.Next() {
			var h dto.HourlyStats
			if err := rowsH.Scan(&h.Hour, &h.Count); err == nil {
				stats.HourlyDistribution = append(stats.HourlyDistribution, h)
			}
		}
	}

	serviceQuery := `
		SELECT s.id, s.name, COUNT(t.id)
		FROM services s
		LEFT JOIN queue_tickets t ON s.id = t.service_id AND t.queue_date = $3 AND t.deleted_at IS NULL
		WHERE s.organization_id = $1 AND ($2::bigint = 0 OR s.branch_id = $2) AND s.deleted_at IS NULL
		GROUP BY s.id, s.name ORDER BY COUNT(t.id) DESC
	`
	rowsS, err := r.db.QueryContext(ctx, serviceQuery, orgID, branchID, today)
	if err == nil {
		defer rowsS.Close()
		for rowsS.Next() {
			var s dto.ServiceStats
			if err := rowsS.Scan(&s.ServiceID, &s.ServiceName, &s.Count); err == nil {
				stats.ServiceDistribution = append(stats.ServiceDistribution, s)
			}
		}
	}

	return stats, nil
}
