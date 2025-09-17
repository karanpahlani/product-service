package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `json:"id" dynamodbav:"id"`
	Name        string    `json:"name" dynamodbav:"name"`
	Description string    `json:"description" dynamodbav:"description"`
	Price       float64   `json:"price" dynamodbav:"price"`
	Category    string    `json:"category" dynamodbav:"category"`
	SKU         string    `json:"sku" dynamodbav:"sku"`
	Stock       int       `json:"stock" dynamodbav:"stock"`
	IsActive    bool      `json:"is_active" dynamodbav:"is_active"`
	CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Stock       int     `json:"stock" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Category    *string  `json:"category,omitempty"`
	SKU         *string  `json:"sku,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

func NewProduct(req CreateProductRequest) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		SKU:         req.SKU,
		Stock:       req.Stock,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Product) Update(req UpdateProductRequest) {
	now := time.Now()
	
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.Category != nil {
		p.Category = *req.Category
	}
	if req.SKU != nil {
		p.SKU = *req.SKU
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}
	
	p.UpdatedAt = now
}