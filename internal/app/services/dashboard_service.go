package services

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/repositories"
)

// DashboardService define el contrato del servicio de dashboard
type DashboardService interface {
	GetStats(companyID uint) (*dtos.DashboardStatsDTO, error)
}

type dashboardService struct {
	dashboardRepo repositories.DashboardRepository
}

// NewDashboardService crea una nueva instancia de DashboardService
func NewDashboardService(dashboardRepo repositories.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetStats(companyID uint) (*dtos.DashboardStatsDTO, error) {
	return s.dashboardRepo.GetStats(companyID)
}
