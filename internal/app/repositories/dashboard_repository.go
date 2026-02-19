package repositories

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/database"
	"time"
)

// DashboardRepository define el contrato del repositorio de dashboard
type DashboardRepository interface {
	GetStats(companyID uint) (*dtos.DashboardStatsDTO, error)
}

// dashboardRepository es la implementación con GORM
type dashboardRepository struct{}

// NewDashboardRepository crea una nueva instancia de DashboardRepository
func NewDashboardRepository() DashboardRepository {
	return &dashboardRepository{}
}

func (r *dashboardRepository) GetStats(companyID uint) (*dtos.DashboardStatsDTO, error) {
	stats := &dtos.DashboardStatsDTO{}

	// ========== JOBS STATS ==========
	// Total jobs
	var totalJobs int64
	if err := database.DB.Table("jobs").Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalJobs).Error; err != nil {
		return nil, err
	}
	stats.TotalJobs = int(totalJobs)

	// Active jobs (published)
	var activeJobs int64
	if err := database.DB.Table("jobs").Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "published").Count(&activeJobs).Error; err != nil {
		return nil, err
	}
	stats.ActiveJobs = int(activeJobs)

	// Draft jobs
	var draftJobs int64
	if err := database.DB.Table("jobs").Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "draft").Count(&draftJobs).Error; err != nil {
		return nil, err
	}
	stats.DraftJobs = int(draftJobs)

	// Closed jobs
	var closedJobs int64
	if err := database.DB.Table("jobs").Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "closed").Count(&closedJobs).Error; err != nil {
		return nil, err
	}
	stats.ClosedJobs = int(closedJobs)

	// ========== CANDIDATES STATS ==========
	var totalCandidates int64
	if err := database.DB.Table("candidates").Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalCandidates).Error; err != nil {
		return nil, err
	}
	stats.TotalCandidates = int(totalCandidates)

	// ========== APPLICATIONS STATS ==========
	var totalApplications int64
	if err := database.DB.Table("applications").Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalApplications).Error; err != nil {
		return nil, err
	}
	stats.TotalApplications = int(totalApplications)

	// Applications by stage
	stats.ApplicationsByStage = dtos.ApplicationsByStageDTO{}

	type stageCount struct {
		Stage string
		Count int
	}
	var stageCounts []stageCount
	if err := database.DB.Table("applications").
		Select("stage, COUNT(*) as count").
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Group("stage").
		Scan(&stageCounts).Error; err != nil {
		return nil, err
	}

	for _, sc := range stageCounts {
		switch sc.Stage {
		case "applied":
			stats.ApplicationsByStage.Applied = sc.Count
		case "screening":
			stats.ApplicationsByStage.Screening = sc.Count
		case "technical":
			stats.ApplicationsByStage.Technical = sc.Count
		case "offer":
			stats.ApplicationsByStage.Offer = sc.Count
		case "hired":
			stats.ApplicationsByStage.Hired = sc.Count
		}
	}

	// ========== MONTHLY STATS ==========
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)

	// New candidates this month
	var newCandidates int64
	if err := database.DB.Table("candidates").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, startOfMonth).
		Count(&newCandidates).Error; err != nil {
		return nil, err
	}
	stats.NewCandidatesThisMonth = int(newCandidates)

	// New applications this month
	var newApplications int64
	if err := database.DB.Table("applications").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, startOfMonth).
		Count(&newApplications).Error; err != nil {
		return nil, err
	}
	stats.NewApplicationsThisMonth = int(newApplications)

	// Hired this month
	var hiredThisMonth int64
	if err := database.DB.Table("applications").
		Where("company_id = ? AND stage = ? AND hired_at >= ? AND deleted_at IS NULL", companyID, "hired", startOfMonth).
		Count(&hiredThisMonth).Error; err != nil {
		return nil, err
	}
	stats.HiredThisMonth = int(hiredThisMonth)

	// ========== PERFORMANCE METRICS ==========
	// Average time to hire (days)
	type avgResult struct {
		AvgDays float64
	}
	var avgTimeToHire avgResult
	if err := database.DB.Table("applications").
		Select("AVG(EXTRACT(EPOCH FROM (hired_at - applied_at)) / 86400) as avg_days").
		Where("company_id = ? AND stage = ? AND hired_at IS NOT NULL AND deleted_at IS NULL", companyID, "hired").
		Scan(&avgTimeToHire).Error; err != nil {
		return nil, err
	}
	stats.AverageTimeToHireDays = avgTimeToHire.AvgDays

	// Conversion rate
	if totalApplications > 0 {
		stats.ConversionRate = float64(stats.ApplicationsByStage.Hired) / float64(totalApplications) * 100
	}

	// ========== TRENDS (Last 30 days) ==========
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	// Applications trend
	type dailyCount struct {
		Date  time.Time
		Count int
	}
	var applicationsTrend []dailyCount
	if err := database.DB.Table("applications").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, thirtyDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&applicationsTrend).Error; err != nil {
		return nil, err
	}

	stats.ApplicationsTrend = make([]dtos.DailyCountDTO, len(applicationsTrend))
	for i, at := range applicationsTrend {
		stats.ApplicationsTrend[i] = dtos.DailyCountDTO{
			Date:  at.Date.Format("2006-01-02"),
			Count: at.Count,
		}
	}

	// Candidates trend
	var candidatesTrend []dailyCount
	if err := database.DB.Table("candidates").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, thirtyDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&candidatesTrend).Error; err != nil {
		return nil, err
	}

	stats.CandidatesTrend = make([]dtos.DailyCountDTO, len(candidatesTrend))
	for i, ct := range candidatesTrend {
		stats.CandidatesTrend[i] = dtos.DailyCountDTO{
			Date:  ct.Date.Format("2006-01-02"),
			Count: ct.Count,
		}
	}

	// ========== TOP JOBS ==========
	type topJob struct {
		ID               uint
		Title            string
		Status           string
		ApplicationCount int
	}
	var topJobs []topJob
	if err := database.DB.Table("jobs").
		Select("jobs.id, jobs.title, jobs.status, COUNT(applications.id) as application_count").
		Joins("LEFT JOIN applications ON applications.job_id = jobs.id AND applications.deleted_at IS NULL").
		Where("jobs.company_id = ? AND jobs.deleted_at IS NULL", companyID).
		Group("jobs.id, jobs.title, jobs.status").
		Order("application_count DESC").
		Limit(5).
		Scan(&topJobs).Error; err != nil {
		return nil, err
	}

	stats.TopJobs = make([]dtos.TopJobDTO, len(topJobs))
	for i, tj := range topJobs {
		stats.TopJobs[i] = dtos.TopJobDTO{
			ID:               tj.ID,
			Title:            tj.Title,
			Status:           tj.Status,
			ApplicationCount: tj.ApplicationCount,
		}
	}

	// ========== CANDIDATE SOURCES ==========
	type sourceCount struct {
		Source string
		Count  int
	}
	var sourceCounts []sourceCount
	if err := database.DB.Table("candidates").
		Select("COALESCE(source, 'unknown') as source, COUNT(*) as count").
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Group("source").
		Order("count DESC").
		Scan(&sourceCounts).Error; err != nil {
		return nil, err
	}

	stats.CandidateSources = make([]dtos.SourceCountDTO, len(sourceCounts))
	for i, sc := range sourceCounts {
		stats.CandidateSources[i] = dtos.SourceCountDTO{
			Source: sc.Source,
			Count:  sc.Count,
		}
	}

	return stats, nil
}
