package repository

import (
	"golang-app/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *models.Review) error
	FindByUserAndProduct(userID, productID uint) (*models.Review, error)
	FindByProductID(productID uint) ([]models.Review, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByUserAndProduct(userID, productID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&review).Error
	return &review, err
}

func (r *reviewRepository) FindByProductID(productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Preload("User").Where("product_id = ?", productID).Find(&reviews).Error
	return reviews, err
}
