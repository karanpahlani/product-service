package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"product-service/internal/models"
	"product-service/internal/service"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	args := m.Called(req)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetProduct(id string) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetAllProducts() ([]*models.Product, error) {
	args := m.Called()
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductService) GetProductsByCategory(category string) ([]*models.Product, error) {
	args := m.Called(category)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(id string, req models.UpdateProductRequest) (*models.Product, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) DeleteProduct(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter(handler *ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	api := router.Group("/api/v1")
	api.GET("/health", handler.HealthCheck)
	
	products := api.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.GET("", handler.GetAllProducts)
		products.GET("/category", handler.GetProductsByCategory)
		products.GET("/:id", handler.GetProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
	}
	
	return router
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	req := models.CreateProductRequest{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "electronics",
		SKU:         "TEST-001",
		Stock:       10,
	}

	product := &models.Product{
		ID:          "test-id",
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		SKU:         req.SKU,
		Stock:       req.Stock,
		IsActive:    true,
	}

	mockService.On("CreateProduct", req).Return(product, nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, product.ID, response.ID)
	assert.Equal(t, product.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_GetProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	product := &models.Product{
		ID:    "test-id",
		Name:  "Test Product",
		Price: 99.99,
	}

	mockService.On("GetProduct", "test-id").Return(product, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/products/test-id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, product.ID, response.ID)

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	mockService.On("GetProduct", "nonexistent-id").Return(nil, service.ErrProductNotFound)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/products/nonexistent-id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetAllProducts_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	products := []*models.Product{
		{ID: "1", Name: "Product 1"},
		{ID: "2", Name: "Product 2"},
	}

	mockService.On("GetAllProducts").Return(products, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/products", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(2), response["count"])

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProductsByCategory_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	products := []*models.Product{
		{ID: "1", Name: "Product 1", Category: "electronics"},
	}

	mockService.On("GetProductsByCategory", "electronics").Return(products, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/products/category?category=electronics", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "electronics", response["category"])
	assert.Equal(t, float64(1), response["count"])

	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProductsByCategory_MissingCategory(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/products/category", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestProductHandler_UpdateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	newName := "Updated Product"
	req := models.UpdateProductRequest{
		Name: &newName,
	}

	updatedProduct := &models.Product{
		ID:   "test-id",
		Name: newName,
	}

	mockService.On("UpdateProduct", "test-id", req).Return(updatedProduct, nil)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/v1/products/test-id", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, updatedProduct.Name, response.Name)

	mockService.AssertExpectations(t)
}

func TestProductHandler_UpdateProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	req := models.UpdateProductRequest{}

	mockService.On("UpdateProduct", "nonexistent-id", req).Return(nil, service.ErrProductNotFound)

	reqBody, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/v1/products/nonexistent-id", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteProduct", "test-id").Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/v1/products/test-id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response["message"], "deleted successfully")

	mockService.AssertExpectations(t)
}

func TestProductHandler_DeleteProduct_NotFound(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	mockService.On("DeleteProduct", "nonexistent-id").Return(service.ErrProductNotFound)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/v1/products/nonexistent-id", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestProductHandler_HealthCheck(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/health", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "product-service", response["service"])
}