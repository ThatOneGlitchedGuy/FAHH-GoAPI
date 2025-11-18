package service

import (
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
)

type StatsService interface {
	GetAdminStats() (*schemas.AdminStatsOut, error)
}

type statsService struct {
	repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository) StatsService {
	return &statsService{repo: repo}
}

func (s *statsService) GetAdminStats() (*schemas.AdminStatsOut, error) {
	totalUsers, err := s.repo.CountUsers()
	if err != nil {
		return nil, err
	}
	totalProducts, err := s.repo.CountProducts()
	if err != nil {
		return nil, err
	}
	totalOrders, err := s.repo.CountOrders()
	if err != nil {
		return nil, err
	}

	return &schemas.AdminStatsOut{
		TotalUsers:   totalUsers,
		TotalProducts: totalProducts,
		TotalOrders:  totalOrders,
	}, nil
}
