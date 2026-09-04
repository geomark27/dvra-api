package repository

import (
	"dvra-api/internal/shared/models"
	"gorm.io/gorm"

	"dvra-api/internal/modules/dashboard/service"
	"time"
)

// DashboardRepository define el contrato del repositorio de dashboard
// dashboardRepository es la implementación con GORM
type DashboardRepository struct{ db *gorm.DB }

// NewDashboardRepository crea una nueva instancia de DashboardRepository
// NewDashboardRepository recibe la conexión inyectada: sin estado global
// (ADR-001 regla 4).
func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetStats(companyID uint) (*service.DashboardStatsDTO, error) {
	stats := &service.DashboardStatsDTO{}

	// ========== JOBS STATS ==========
	// Total jobs
	var totalJobs int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalJobs).Error; err != nil {
		return nil, err
	}
	stats.TotalJobs = int(totalJobs)

	// Active jobs (published)
	var activeJobs int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "published").Count(&activeJobs).Error; err != nil {
		return nil, err
	}
	stats.ActiveJobs = int(activeJobs)

	// Draft jobs
	var draftJobs int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "draft").Count(&draftJobs).Error; err != nil {
		return nil, err
	}
	stats.DraftJobs = int(draftJobs)

	// Closed jobs
	var closedJobs int64
	if err := r.db.Model(&models.Job{}).Where("company_id = ? AND status = ? AND deleted_at IS NULL", companyID, "closed").Count(&closedJobs).Error; err != nil {
		return nil, err
	}
	stats.ClosedJobs = int(closedJobs)

	// ========== CANDIDATES STATS ==========
	var totalCandidates int64
	if err := r.db.Model(&models.Candidate{}).Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalCandidates).Error; err != nil {
		return nil, err
	}
	stats.TotalCandidates = int(totalCandidates)

	// ========== APPLICATIONS STATS ==========
	var totalApplications int64
	if err := r.db.Model(&models.Application{}).Where("company_id = ? AND deleted_at IS NULL", companyID).Count(&totalApplications).Error; err != nil {
		return nil, err
	}
	stats.TotalApplications = int(totalApplications)

	// Applications by stage
	stats.ApplicationsByStage = service.ApplicationsByStageDTO{}

	type stageCount struct {
		Stage string
		Count int
	}
	var stageCounts []stageCount
	if err := r.db.Model(&models.Application{}).
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
	if err := r.db.Model(&models.Candidate{}).
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, startOfMonth).
		Count(&newCandidates).Error; err != nil {
		return nil, err
	}
	stats.NewCandidatesThisMonth = int(newCandidates)

	// New applications this month
	var newApplications int64
	if err := r.db.Model(&models.Application{}).
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, startOfMonth).
		Count(&newApplications).Error; err != nil {
		return nil, err
	}
	stats.NewApplicationsThisMonth = int(newApplications)

	// Hired this month
	var hiredThisMonth int64
	if err := r.db.Model(&models.Application{}).
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
	if err := r.db.Model(&models.Application{}).
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
	if err := r.db.Model(&models.Application{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, thirtyDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&applicationsTrend).Error; err != nil {
		return nil, err
	}

	stats.ApplicationsTrend = make([]service.DailyCountDTO, len(applicationsTrend))
	for i, at := range applicationsTrend {
		stats.ApplicationsTrend[i] = service.DailyCountDTO{
			Date:  at.Date.Format("2006-01-02"),
			Count: at.Count,
		}
	}

	// Candidates trend
	var candidatesTrend []dailyCount
	if err := r.db.Model(&models.Candidate{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("company_id = ? AND created_at >= ? AND deleted_at IS NULL", companyID, thirtyDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&candidatesTrend).Error; err != nil {
		return nil, err
	}

	stats.CandidatesTrend = make([]service.DailyCountDTO, len(candidatesTrend))
	for i, ct := range candidatesTrend {
		stats.CandidatesTrend[i] = service.DailyCountDTO{
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
	if err := r.db.Model(&models.Job{}).
		Select("jobs.id, jobs.title, jobs.status, COUNT(applications.id) as application_count").
		Joins("LEFT JOIN applications ON applications.job_id = jobs.id AND applications.deleted_at IS NULL").
		Where("jobs.company_id = ? AND jobs.deleted_at IS NULL", companyID).
		Group("jobs.id, jobs.title, jobs.status").
		Order("application_count DESC").
		Limit(5).
		Scan(&topJobs).Error; err != nil {
		return nil, err
	}

	stats.TopJobs = make([]service.TopJobDTO, len(topJobs))
	for i, tj := range topJobs {
		stats.TopJobs[i] = service.TopJobDTO{
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
	if err := r.db.Model(&models.Candidate{}).
		Select("COALESCE(source, 'unknown') as source, COUNT(*) as count").
		Where("company_id = ? AND deleted_at IS NULL", companyID).
		Group("source").
		Order("count DESC").
		Scan(&sourceCounts).Error; err != nil {
		return nil, err
	}

	stats.CandidateSources = make([]service.SourceCountDTO, len(sourceCounts))
	for i, sc := range sourceCounts {
		stats.CandidateSources[i] = service.SourceCountDTO{
			Source: sc.Source,
			Count:  sc.Count,
		}
	}

	return stats, nil
}
