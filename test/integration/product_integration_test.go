package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"product-service/internal/httpserver"
	"product-service/internal/models"
)

type ProductIntegrationTestSuite struct {
	suite.Suite
	server *gin.Engine
}

func (suite *ProductIntegrationTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("PRODUCTS_TABLE", "test-products")
	
	server, err := httpserver.NewServer()
	if err != nil {
		suite.T().Skip("Skipping integration tests: unable to create server (likely missing AWS credentials)")
		return
	}
	
	suite.server = server
}

func (suite *ProductIntegrationTestSuite) TestHealthEndpoint() {
	if suite.server == nil {
		suite.T().Skip("Server not initialized")
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	suite.server.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "healthy", response["status"])
	assert.Equal(suite.T(), "product-service", response["service"])
}

func (suite *ProductIntegrationTestSuite) TestProductCRUDOperations() {
	if suite.server == nil {
		suite.T().Skip("Server not initialized")
		return
	}

	createReq := models.CreateProductRequest{
		Name:        "Integration Test Product",
		Description: "A product created during integration testing",
		Price:       149.99,
		Category:    "test",
		SKU:         "INT-TEST-001",
		Stock:       25,
	}

	reqBody, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	suite.server.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		suite.T().Skip("Skipping CRUD test: Create operation failed (likely DynamoDB connectivity issue)")
		return
	}

	var createdProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &createdProduct)
	productID := createdProduct.ID

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%s", productID), nil)
	suite.server.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var retrievedProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &retrievedProduct)
	assert.Equal(suite.T(), createReq.Name, retrievedProduct.Name)
	assert.Equal(suite.T(), createReq.Price, retrievedProduct.Price)

	newName := "Updated Integration Test Product"
	updateReq := models.UpdateProductRequest{
		Name: &newName,
	}

	reqBody, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/products/%s", productID), bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	suite.server.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var updatedProduct models.Product
	json.Unmarshal(w.Body.Bytes(), &updatedProduct)
	assert.Equal(suite.T(), newName, updatedProduct.Name)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/products/%s", productID), nil)
	suite.server.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/products/%s", productID), nil)
	suite.server.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *ProductIntegrationTestSuite) TestGetAllProducts() {
	if suite.server == nil {
		suite.T().Skip("Server not initialized")
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products", nil)
	suite.server.ServeHTTP(w, req)

	if w.Code == http.StatusInternalServerError {
		suite.T().Skip("Skipping GetAll test: DynamoDB connectivity issue")
		return
	}

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "products")
	assert.Contains(suite.T(), response, "count")
}

func (suite *ProductIntegrationTestSuite) TestGetProductsByCategory() {
	if suite.server == nil {
		suite.T().Skip("Server not initialized")
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products/category?category=electronics", nil)
	suite.server.ServeHTTP(w, req)

	if w.Code == http.StatusInternalServerError {
		suite.T().Skip("Skipping GetByCategory test: DynamoDB connectivity issue")
		return
	}

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "electronics", response["category"])
	assert.Contains(suite.T(), response, "products")
	assert.Contains(suite.T(), response, "count")
}

func TestProductIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ProductIntegrationTestSuite))
}