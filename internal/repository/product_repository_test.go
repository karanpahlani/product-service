package repository

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"product-service/internal/database"
	"product-service/internal/models"
)

type MockDynamoDBClient struct {
	mock.Mock
}

func (m *MockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *MockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.ScanOutput), args.Error(1)
}

func (m *MockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.DeleteItemOutput), args.Error(1)
}

func createTestProduct() *models.Product {
	return &models.Product{
		ID:          "test-id",
		Name:        "Test Product",
		Description: "A test product",
		Price:       99.99,
		Category:    "electronics",
		SKU:         "TEST-001",
		Stock:       10,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func TestProductRepository_Create_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	product := createTestProduct()

	mockClient.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := repo.Create(product)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_GetByID_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	product := createTestProduct()
	item, _ := dynamodbattribute.MarshalMap(product)

	output := &dynamodb.GetItemOutput{
		Item: item,
	}

	mockClient.On("GetItem", mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
		return *input.TableName == "test-table" && 
			   *input.Key["id"].S == "test-id"
	})).Return(output, nil)

	result, err := repo.GetByID("test-id")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, product.ID, result.ID)
	assert.Equal(t, product.Name, result.Name)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	output := &dynamodb.GetItemOutput{
		Item: nil,
	}

	mockClient.On("GetItem", mock.AnythingOfType("*dynamodb.GetItemInput")).Return(output, nil)

	result, err := repo.GetByID("nonexistent-id")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_GetAll_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	product1 := createTestProduct()
	product1.ID = "id-1"
	product2 := createTestProduct()
	product2.ID = "id-2"

	item1, _ := dynamodbattribute.MarshalMap(product1)
	item2, _ := dynamodbattribute.MarshalMap(product2)

	output := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			item1,
			item2,
		},
	}

	mockClient.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
		return *input.TableName == "test-table" &&
			   input.FilterExpression != nil &&
			   *input.FilterExpression == "is_active = :active"
	})).Return(output, nil)

	results, err := repo.GetAll()

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "id-1", results[0].ID)
	assert.Equal(t, "id-2", results[1].ID)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_GetByCategory_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	product := createTestProduct()
	item, _ := dynamodbattribute.MarshalMap(product)

	output := &dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			item,
		},
	}

	mockClient.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
		return *input.TableName == "test-table" &&
			   input.FilterExpression != nil &&
			   *input.FilterExpression == "category = :category AND is_active = :active" &&
			   *input.ExpressionAttributeValues[":category"].S == "electronics"
	})).Return(output, nil)

	results, err := repo.GetByCategory("electronics")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "electronics", results[0].Category)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_Update_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	product := createTestProduct()

	mockClient.On("PutItem", mock.AnythingOfType("*dynamodb.PutItemInput")).Return(&dynamodb.PutItemOutput{}, nil)

	err := repo.Update(product)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestProductRepository_Delete_Success(t *testing.T) {
	mockClient := new(MockDynamoDBClient)
	db := &database.DynamoDBClient{
		Client:    mockClient,
		TableName: "test-table",
	}
	repo := NewProductRepository(db)

	mockClient.On("DeleteItem", mock.MatchedBy(func(input *dynamodb.DeleteItemInput) bool {
		return *input.TableName == "test-table" &&
			   *input.Key["id"].S == "test-id"
	})).Return(&dynamodb.DeleteItemOutput{}, nil)

	err := repo.Delete("test-id")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}