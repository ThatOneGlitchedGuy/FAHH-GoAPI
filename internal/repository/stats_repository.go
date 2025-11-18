package repository

import (
	"golang-app/internal/models"
	"gorm.io/gorm"
)

type StatsRepository interface {
	CountUsers() (int64, error)
	CountProducts() (int64, error)
	CountOrders() (int64, error)
}

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}

func (r *statsRepository) CountProducts() (int64, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Count(&count).Error
	return count, err
}

func (r *statsRepository) CountOrders() (int64, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Count(&count).Error
	return count, err
}
