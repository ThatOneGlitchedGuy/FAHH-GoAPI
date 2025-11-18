package repository

import (
	"golang-app/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(id uint) (*models.Product, error)
	FindActiveProducts(offset, limit int) ([]models.Product, int64, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(product *models.Product) error
	FindProductByAgent(productID, agentID uint) (*models.Product, error)
	FindProductsByIDs(productIDs []uint) ([]models.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *productRepository) FindActiveProducts(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	r.db.Model(&models.Product{}).Where("is_active = ?", true).Count(&total)

	err := r.db.Where("is_active = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error
	return products, total, err
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(product *models.Product) error {
	return r.db.Delete(product).Error
}

func (r *productRepository) FindProductByAgent(productID, agentID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ? AND agent_id = ?", productID, agentID).First(&product).Error
	return &product, err
}

func (r *productRepository) FindProductsByIDs(productIDs []uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("id IN (?) AND is_active = ?", productIDs, true).Find(&products).Error
	return products, err
}
