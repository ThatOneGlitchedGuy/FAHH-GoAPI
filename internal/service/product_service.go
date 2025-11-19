package service

import (
	"errors"
	"fmt"
	"golang-app/internal/config"
	"golang-app/internal/models"
	"golang-app/internal/repository"
	"golang-app/internal/schemas"
	"golang-app/internal/security"
	"math"
	"strings"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type ProductService interface {
	ListProducts(page, size int) (*schemas.Page, error)
	GetProductByID(id uint) (*schemas.ProductOut, error)
	CreateProduct(data *schemas.ProductCreate, agentID uint) (*schemas.ProductOut, error)
	UpdateProduct(productID uint, data *schemas.ProductUpdate, agentID uint) (*schemas.ProductOut, error)
	DeleteProduct(productID uint, agentID uint) error
	CreateReview(productID, userID uint, data *schemas.ReviewCreate) (*schemas.ReviewOut, error)
	ListReviews(productID uint) ([]schemas.ReviewOut, error)
	HasUserPurchasedProduct(userID, productID uint) (bool, error)
}

type productService struct {
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	reviewRepo  repository.ReviewRepository
	userRepo    repository.UserRepository
	config      *config.Settings
}

func NewProductService(
	productRepo repository.ProductRepository,
	orderRepo repository.OrderRepository,
	reviewRepo repository.ReviewRepository,
	userRepo repository.UserRepository,
	cfg *config.Settings,
) ProductService {
	return &productService{
		productRepo: productRepo,
		orderRepo:   orderRepo,
		reviewRepo:  reviewRepo,
		userRepo:    userRepo,
		config:      cfg,
	}
}

func (s *productService) productToProductOut(product *models.Product) *schemas.ProductOut {
	desc := product.Description
	return &schemas.ProductOut{
		ID:          product.ID,
		AgentID:     product.AgentID,
		Title:       product.Title,
		Description: &desc,
		Price:       product.Price,
		IsActive:    product.IsActive,
		CreatedAt:   product.CreatedAt,
	}
}

func (s *productService) ListProducts(page, size int) (*schemas.Page, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = s.config.PageSizeDefault
	}
	if size > s.config.PageSizeMax {
		size = s.config.PageSizeMax
	}

	cacheKey := fmt.Sprintf("products:page:%d:size:%d", page, size)
	if cached, found := Cache.Get(cacheKey); found {
		return cached.(*schemas.Page), nil
	}

	offset := (page - 1) * size
	products, total, err := s.productRepo.FindActiveProducts(offset, size)
	if err != nil {
		return nil, err
	}

	productOuts := make([]schemas.ProductOut, len(products))
	for i, product := range products {
		productOuts[i] = *s.productToProductOut(&product)
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	meta := schemas.PageMeta{
		Page:  page,
		Size:  size,
		Total: total,
	}

	pageResult := &schemas.Page{
		Meta:  meta,
		Items: productOuts,
	}

	Cache.Set(cacheKey, pageResult, cache.DefaultExpiration)

	return pageResult, nil
}

func (s *productService) GetProductByID(id uint) (*schemas.ProductOut, error) {
	cacheKey := fmt.Sprintf("product:%d", id)
	if cached, found := Cache.Get(cacheKey); found {
		return cached.(*schemas.ProductOut), nil
	}

	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if !product.IsActive {
		return nil, gorm.ErrRecordNotFound
	}

	productOut := s.productToProductOut(product)
	Cache.Set(cacheKey, productOut, cache.DefaultExpiration)

	return productOut, nil
}

func (s *productService) CreateProduct(data *schemas.ProductCreate, agentID uint) (*schemas.ProductOut, error) {
	product := &models.Product{
		AgentID:     agentID,
		Title:       data.Title,
		Description: *data.Description,
		Price:       data.Price,
		IsActive:    data.IsActive,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	// Invalidate list cache
	for key := range Cache.Items() {
		if strings.HasPrefix(key, "products:page:") {
			Cache.Delete(key)
		}
	}

	return s.productToProductOut(product), nil
}

func (s *productService) UpdateProduct(productID uint, data *schemas.ProductUpdate, agentID uint) (*schemas.ProductOut, error) {
	product, err := s.productRepo.FindProductByAgent(productID, agentID)
	if err != nil {
		return nil, err
	}

	if data.Title != nil {
		product.Title = *data.Title
	}
	if data.Description != nil {
		product.Description = *data.Description
	}
	if data.Price != nil {
		product.Price = *data.Price
	}
	if data.IsActive != nil {
		product.IsActive = *data.IsActive
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("product:%d", productID)
	Cache.Delete(cacheKey)
	for key := range Cache.Items() {
		if strings.HasPrefix(key, "products:page:") {
			Cache.Delete(key)
		}
	}

	return s.productToProductOut(product), nil
}

func (s *productService) DeleteProduct(productID uint, agentID uint) error {
	product, err := s.productRepo.FindProductByAgent(productID, agentID)
	if err != nil {
		return err
	}

	if err := s.productRepo.Delete(product); err != nil {
		return err
	}

	// Invalidate caches
	cacheKey := fmt.Sprintf("product:%d", productID)
	Cache.Delete(cacheKey)
	for key := range Cache.Items() {
		if strings.HasPrefix(key, "products:page:") {
			Cache.Delete(key)
		}
	}

	return nil
}

func (s *productService) HasUserPurchasedProduct(userID, productID uint) (bool, error) {
	return s.orderRepo.HasUserPurchasedProduct(userID, productID)
}

func (s *productService) CreateReview(productID, userID uint, data *schemas.ReviewCreate) (*schemas.ReviewOut, error) {
	purchased, err := s.HasUserPurchasedProduct(userID, productID)
	if err != nil {
		return nil, err
	}
	if !purchased {
		return nil, errors.New("user has not purchased this product")
	}

	existingReview, err := s.reviewRepo.FindByUserAndProduct(userID, productID)
	if err == nil && existingReview != nil {
		return nil, errors.New("user has already reviewed this product")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	review := &models.Review{
		ProductID: productID,
		UserID:    userID,
		Rating:    data.Rating,
		Comment:   data.Comment,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var fullName string
	if user.FullNameEncrypted != nil {
		decrypted, err := security.Decrypt(user.FullNameEncrypted, s.config.FernetKey, 0)
		if err == nil {
			fullName = string(decrypted)
		}
	}

	return &schemas.ReviewOut{
		ID:        review.ID,
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		User: &schemas.UserReviewerOut{
			ID:       user.ID,
			FullName: fullName,
		},
	}, nil
}

func (s *productService) ListReviews(productID uint) ([]schemas.ReviewOut, error) {
	reviews, err := s.reviewRepo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	reviewOuts := make([]schemas.ReviewOut, len(reviews))
	for i, review := range reviews {
		var fullName string
		if review.User.FullNameEncrypted != nil {
			decrypted, err := security.Decrypt(review.User.FullNameEncrypted, s.config.FernetKey, 0)
			if err == nil {
				fullName = string(decrypted)
			}
		}

		reviewOuts[i] = schemas.ReviewOut{
			ID:        review.ID,
			ProductID: review.ProductID,
			UserID:    review.UserID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
			User: &schemas.UserReviewerOut{
				ID:       review.User.ID,
				FullName: fullName,
			},
		}
	}
	return reviewOuts, nil
}
