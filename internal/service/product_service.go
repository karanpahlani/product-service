package service

import (
	"errors"
	"fmt"

	"product-service/internal/models"
	"product-service/internal/repository"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidProduct  = errors.New("invalid product data")
)

type ProductService interface {
	CreateProduct(req models.CreateProductRequest) (*models.Product, error)
	GetProduct(id string) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	GetProductsByCategory(category string) ([]*models.Product, error)
	UpdateProduct(id string, req models.UpdateProductRequest) (*models.Product, error)
	DeleteProduct(id string) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidProduct, err)
	}

	product := models.NewProduct(req)

	if err := s.repo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProduct(id string) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: product ID cannot be empty", ErrInvalidProduct)
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	return product, nil
}

func (s *productService) GetAllProducts() ([]*models.Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return products, nil
}

func (s *productService) GetProductsByCategory(category string) ([]*models.Product, error) {
	if category == "" {
		return nil, fmt.Errorf("%w: category cannot be empty", ErrInvalidProduct)
	}

	products, err := s.repo.GetByCategory(category)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category: %w", err)
	}

	return products, nil
}

func (s *productService) UpdateProduct(id string, req models.UpdateProductRequest) (*models.Product, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: product ID cannot be empty", ErrInvalidProduct)
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product for update: %w", err)
	}

	if product == nil {
		return nil, ErrProductNotFound
	}

	if err := s.validateUpdateRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidProduct, err)
	}

	product.Update(req)

	if err := s.repo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) DeleteProduct(id string) error {
	if id == "" {
		return fmt.Errorf("%w: product ID cannot be empty", ErrInvalidProduct)
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get product for deletion: %w", err)
	}

	if product == nil {
		return ErrProductNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (s *productService) validateCreateRequest(req models.CreateProductRequest) error {
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if req.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if req.Category == "" {
		return errors.New("product category is required")
	}
	if req.SKU == "" {
		return errors.New("product SKU is required")
	}
	if req.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	return nil
}

func (s *productService) validateUpdateRequest(req models.UpdateProductRequest) error {
	if req.Price != nil && *req.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if req.Stock != nil && *req.Stock < 0 {
		return errors.New("product stock cannot be negative")
	}
	if req.Name != nil && *req.Name == "" {
		return errors.New("product name cannot be empty")
	}
	if req.Category != nil && *req.Category == "" {
		return errors.New("product category cannot be empty")
	}
	if req.SKU != nil && *req.SKU == "" {
		return errors.New("product SKU cannot be empty")
	}
	return nil
}