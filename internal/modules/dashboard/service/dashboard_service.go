package service

import ()

// DashboardService define el contrato del servicio de dashboard
type DashboardService interface {
	GetStats(companyID uint) (*DashboardStatsDTO, error)
}

type dashboardService struct {
	dashboardRepo DashboardRepository
}

// NewDashboardService crea una nueva instancia de DashboardService
func NewDashboardService(dashboardRepo DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetStats(companyID uint) (*DashboardStatsDTO, error) {
	return s.dashboardRepo.GetStats(companyID)
}
