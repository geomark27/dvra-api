package dtos

// ApplicationsByStageDTO represents applications count by stage
type ApplicationsByStageDTO struct {
	Applied   int `json:"applied"`
	Screening int `json:"screening"`
	Technical int `json:"technical"`
	Offer     int `json:"offer"`
	Hired     int `json:"hired"`
}

// DailyCountDTO represents a count for a specific date
type DailyCountDTO struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// TopJobDTO represents a job with its application count
type TopJobDTO struct {
	ID               uint   `json:"id"`
	Title            string `json:"title"`
	ApplicationCount int    `json:"application_count"`
	Status           string `json:"status"`
}

// SourceCountDTO represents candidate source statistics
type SourceCountDTO struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

// DashboardStatsDTO represents all dashboard statistics
type DashboardStatsDTO struct {
	// Jobs stats
	TotalJobs  int `json:"total_jobs"`
	ActiveJobs int `json:"active_jobs"`
	DraftJobs  int `json:"draft_jobs"`
	ClosedJobs int `json:"closed_jobs"`

	// Candidates stats
	TotalCandidates int `json:"total_candidates"`

	// Applications stats
	TotalApplications   int                    `json:"total_applications"`
	ApplicationsByStage ApplicationsByStageDTO `json:"applications_by_stage"`

	// Monthly stats
	NewCandidatesThisMonth   int `json:"new_candidates_this_month"`
	NewApplicationsThisMonth int `json:"new_applications_this_month"`
	HiredThisMonth           int `json:"hired_this_month"`

	// Performance metrics
	AverageTimeToHireDays float64 `json:"average_time_to_hire_days"`
	ConversionRate        float64 `json:"conversion_rate"` // Hired / Total Applications

	// Trends (last 30 days)
	ApplicationsTrend []DailyCountDTO `json:"applications_trend"`
	CandidatesTrend   []DailyCountDTO `json:"candidates_trend"`

	// Top jobs by applications
	TopJobs []TopJobDTO `json:"top_jobs"`

	// Candidate sources
	CandidateSources []SourceCountDTO `json:"candidate_sources"`
}
