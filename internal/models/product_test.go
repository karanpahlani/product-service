package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	req := CreateProductRequest{
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "electronics",
		SKU:         "TEST-001",
		Stock:       10,
	}

	product := NewProduct(req)

	assert.NotEmpty(t, product.ID)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Description, product.Description)
	assert.Equal(t, req.Price, product.Price)
	assert.Equal(t, req.Category, product.Category)
	assert.Equal(t, req.SKU, product.SKU)
	assert.Equal(t, req.Stock, product.Stock)
	assert.True(t, product.IsActive)
	assert.False(t, product.CreatedAt.IsZero())
	assert.False(t, product.UpdatedAt.IsZero())
}

func TestProduct_Update(t *testing.T) {
	product := &Product{
		ID:          "test-id",
		Name:        "Original Name",
		Description: "Original Description",
		Price:       50.00,
		Category:    "original",
		SKU:         "ORIG-001",
		Stock:       5,
		IsActive:    true,
		CreatedAt:   time.Now().Add(-time.Hour),
		UpdatedAt:   time.Now().Add(-time.Hour),
	}

	originalUpdatedAt := product.UpdatedAt

	newName := "Updated Name"
	newPrice := 75.00
	newStock := 15
	isActive := false

	updateReq := UpdateProductRequest{
		Name:     &newName,
		Price:    &newPrice,
		Stock:    &newStock,
		IsActive: &isActive,
	}

	product.Update(updateReq)

	assert.Equal(t, newName, product.Name)
	assert.Equal(t, "Original Description", product.Description)
	assert.Equal(t, newPrice, product.Price)
	assert.Equal(t, "original", product.Category)
	assert.Equal(t, "ORIG-001", product.SKU)
	assert.Equal(t, newStock, product.Stock)
	assert.Equal(t, isActive, product.IsActive)
	assert.True(t, product.UpdatedAt.After(originalUpdatedAt))
}

func TestProduct_UpdateWithNilValues(t *testing.T) {
	product := &Product{
		ID:          "test-id",
		Name:        "Original Name",
		Description: "Original Description",
		Price:       50.00,
		Category:    "original",
		SKU:         "ORIG-001",
		Stock:       5,
		IsActive:    true,
		CreatedAt:   time.Now().Add(-time.Hour),
		UpdatedAt:   time.Now().Add(-time.Hour),
	}

	originalValues := *product
	updateReq := UpdateProductRequest{}

	product.Update(updateReq)

	assert.Equal(t, originalValues.Name, product.Name)
	assert.Equal(t, originalValues.Description, product.Description)
	assert.Equal(t, originalValues.Price, product.Price)
	assert.Equal(t, originalValues.Category, product.Category)
	assert.Equal(t, originalValues.SKU, product.SKU)
	assert.Equal(t, originalValues.Stock, product.Stock)
	assert.Equal(t, originalValues.IsActive, product.IsActive)
	assert.True(t, product.UpdatedAt.After(originalValues.UpdatedAt))
}