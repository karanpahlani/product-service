package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"product-service/internal/models"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id string) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll() ([]*models.Product, error) {
	args := m.Called()
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByCategory(category string) ([]*models.Product, error) {
	args := m.Called(category)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductService_CreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := models.CreateProductRequest{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "electronics",
		SKU:         "TEST-001",
		Stock:       10,
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.CreateProduct(req)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Price, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_ValidationError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := models.CreateProductRequest{
		Name:     "",
		Price:    99.99,
		Category: "electronics",
		SKU:      "TEST-001",
		Stock:    10,
	}

	product, err := service.CreateProduct(req)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "invalid product data")
}

func TestProductService_GetProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProduct := &models.Product{
		ID:    "test-id",
		Name:  "Test Product",
		Price: 99.99,
	}

	mockRepo.On("GetByID", "test-id").Return(expectedProduct, nil)

	product, err := service.GetProduct("test-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("GetByID", "nonexistent-id").Return((*models.Product)(nil), nil)

	product, err := service.GetProduct("nonexistent-id")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProduct_EmptyID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	product, err := service.GetProduct("")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "invalid product data")
}

func TestProductService_GetAllProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []*models.Product{
		{ID: "1", Name: "Product 1"},
		{ID: "2", Name: "Product 2"},
	}

	mockRepo.On("GetAll").Return(expectedProducts, nil)

	products, err := service.GetAllProducts()

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductsByCategory_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	expectedProducts := []*models.Product{
		{ID: "1", Name: "Product 1", Category: "electronics"},
	}

	mockRepo.On("GetByCategory", "electronics").Return(expectedProducts, nil)

	products, err := service.GetProductsByCategory("electronics")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	existingProduct := &models.Product{
		ID:    "test-id",
		Name:  "Original Name",
		Price: 50.00,
	}

	newName := "Updated Name"
	updateReq := models.UpdateProductRequest{
		Name: &newName,
	}

	mockRepo.On("GetByID", "test-id").Return(existingProduct, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Product")).Return(nil)

	product, err := service.UpdateProduct("test-id", updateReq)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, newName, product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	updateReq := models.UpdateProductRequest{}

	mockRepo.On("GetByID", "nonexistent-id").Return((*models.Product)(nil), nil)

	product, err := service.UpdateProduct("nonexistent-id", updateReq)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	existingProduct := &models.Product{
		ID:   "test-id",
		Name: "Test Product",
	}

	mockRepo.On("GetByID", "test-id").Return(existingProduct, nil)
	mockRepo.On("Delete", "test-id").Return(nil)

	err := service.DeleteProduct("test-id")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_DeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	mockRepo.On("GetByID", "nonexistent-id").Return((*models.Product)(nil), nil)

	err := service.DeleteProduct("nonexistent-id")

	assert.Error(t, err)
	assert.Equal(t, ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_validateCreateRequest(t *testing.T) {
	service := &productService{}

	tests := []struct {
		name    string
		req     models.CreateProductRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: models.CreateProductRequest{
				Name:     "Test Product",
				Price:    99.99,
				Category: "electronics",
				SKU:      "TEST-001",
				Stock:    10,
			},
			wantErr: false,
		},
		{
			name: "empty name",
			req: models.CreateProductRequest{
				Name:     "",
				Price:    99.99,
				Category: "electronics",
				SKU:      "TEST-001",
				Stock:    10,
			},
			wantErr: true,
			errMsg:  "product name is required",
		},
		{
			name: "zero price",
			req: models.CreateProductRequest{
				Name:     "Test Product",
				Price:    0,
				Category: "electronics",
				SKU:      "TEST-001",
				Stock:    10,
			},
			wantErr: true,
			errMsg:  "product price must be greater than 0",
		},
		{
			name: "negative stock",
			req: models.CreateProductRequest{
				Name:     "Test Product",
				Price:    99.99,
				Category: "electronics",
				SKU:      "TEST-001",
				Stock:    -1,
			},
			wantErr: true,
			errMsg:  "product stock cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateRequest(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}