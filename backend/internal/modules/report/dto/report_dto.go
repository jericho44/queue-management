package dto

type DashboardStatsResponse struct {
	TotalTicketsToday   int64          `json:"total_tickets_today"`
	CompletedToday      int64          `json:"completed_today"`
	WaitingCount        int64          `json:"waiting_count"`
	ServingCount        int64          `json:"serving_count"`
	AvgWaitTimeSec      int64          `json:"avg_wait_time_sec"`
	AvgServiceTimeSec   int64          `json:"avg_service_time_sec"`
	NoShowCount         int64          `json:"no_show_count"`
	CancelledCount      int64          `json:"cancelled_count"`
	ActiveCounters      int64          `json:"active_counters"`
	HourlyDistribution  []HourlyStats  `json:"hourly_distribution"`
	ServiceDistribution []ServiceStats `json:"service_distribution"`
}

type HourlyStats struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type ServiceStats struct {
	ServiceID   int64  `json:"service_id"`
	ServiceName string `json:"service_name"`
	Count       int64  `json:"count"`
}
