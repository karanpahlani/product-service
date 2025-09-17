package repository

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"product-service/internal/database"
	"product-service/internal/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	GetByCategory(category string) ([]*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
}

type productRepository struct {
	db *database.DynamoDBClient
}

func NewProductRepository(db *database.DynamoDBClient) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *models.Product) error {
	item, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.db.TableName),
		Item:      item,
	}

	_, err = r.db.Client.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(id string) (*models.Product, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	result, err := r.db.Client.GetItem(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var product models.Product
	err = dynamodbattribute.UnmarshalMap(result.Item, &product)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

func (r *productRepository) GetAll() ([]*models.Product, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.db.TableName),
		FilterExpression: aws.String("is_active = :active"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":active": {
				BOOL: aws.Bool(true),
			},
		},
	}

	result, err := r.db.Client.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products: %w", err)
	}

	var products []*models.Product
	for _, item := range result.Items {
		var product models.Product
		err = dynamodbattribute.UnmarshalMap(item, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %w", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *productRepository) GetByCategory(category string) ([]*models.Product, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.db.TableName),
		FilterExpression: aws.String("category = :category AND is_active = :active"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":category": {
				S: aws.String(category),
			},
			":active": {
				BOOL: aws.Bool(true),
			},
		},
	}

	result, err := r.db.Client.Scan(input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan products by category: %w", err)
	}

	var products []*models.Product
	for _, item := range result.Items {
		var product models.Product
		err = dynamodbattribute.UnmarshalMap(item, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %w", err)
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *productRepository) Update(product *models.Product) error {
	item, err := dynamodbattribute.MarshalMap(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.db.TableName),
		Item:      item,
	}

	_, err = r.db.Client.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(id string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(r.db.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	}

	_, err := r.db.Client.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}