package product

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateSKU = errors.New("sku already exists")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type CreateProductInput struct {
	Name        *string
	SKU         *string
	Description *string
	Category    *string
	Brand       *string
	Price       *float64
	Stock       *int64
	ImageURL    *string
	IsActive    *bool
}

type UpdateProductInput struct {
	Name        *string
	SKU         *string
	Description *string
	Category    *string
	Brand       *string
	Price       *float64
	Stock       *int64
	ImageURL    *string
	IsActive    *bool
}

type ProductService struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, in CreateProductInput) (*Product, error) {
	normalized, err := normalizeAndValidateCreateInput(in)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	p := &Product{
		Name:        normalized.Name,
		SKU:         normalized.SKU,
		Description: normalized.Description,
		Category:    normalized.Category,
		Brand:       normalized.Brand,
		Price:       normalized.Price,
		Stock:       normalized.Stock,
		ImageURL:    normalized.ImageURL,
		IsActive:    normalized.IsActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.Create(ctx, p)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) Update(ctx context.Context, id int64, in UpdateProductInput) (*Product, error) {
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	normalized, err := normalizeAndValidateUpdateInput(in, current.IsActive)
	if err != nil {
		return nil, err
	}

	p := &Product{
		ID:          id,
		Name:        normalized.Name,
		SKU:         normalized.SKU,
		Description: normalized.Description,
		Category:    normalized.Category,
		Brand:       normalized.Brand,
		Price:       normalized.Price,
		Stock:       normalized.Stock,
		ImageURL:    normalized.ImageURL,
		IsActive:    normalized.IsActive,
		UpdatedAt:   time.Now().UTC(),
	}
	return s.repo.Update(ctx, p)
}

type normalizedProductInput struct {
	Name        string
	SKU         string
	Description string
	Category    string
	Brand       string
	Price       float64
	Stock       int64
	ImageURL    string
	IsActive    bool
}

func normalizeAndValidateCreateInput(in CreateProductInput) (normalizedProductInput, error) {
	return normalizeAndValidateProductInput(
		in.Name,
		in.SKU,
		in.Description,
		in.Category,
		in.Brand,
		in.Price,
		in.Stock,
		in.ImageURL,
		in.IsActive,
		true,
	)
}

func normalizeAndValidateUpdateInput(in UpdateProductInput, currentIsActive bool) (normalizedProductInput, error) {
	return normalizeAndValidateProductInput(
		in.Name,
		in.SKU,
		in.Description,
		in.Category,
		in.Brand,
		in.Price,
		in.Stock,
		in.ImageURL,
		in.IsActive,
		currentIsActive,
	)
}

func normalizeAndValidateProductInput(
	name *string,
	sku *string,
	description *string,
	category *string,
	brand *string,
	price *float64,
	stock *int64,
	imageURL *string,
	isActive *bool,
	defaultIsActive bool,
) (normalizedProductInput, error) {
	out := normalizedProductInput{}

	switch {
	case name == nil:
		return out, &ValidationError{Field: "name", Message: "is required"}
	case strings.TrimSpace(*name) == "":
		return out, &ValidationError{Field: "name", Message: "is required"}
	case sku == nil:
		return out, &ValidationError{Field: "sku", Message: "is required"}
	case strings.TrimSpace(*sku) == "":
		return out, &ValidationError{Field: "sku", Message: "is required"}
	case description == nil:
		return out, &ValidationError{Field: "description", Message: "is required"}
	case strings.TrimSpace(*description) == "":
		return out, &ValidationError{Field: "description", Message: "is required"}
	case category == nil:
		return out, &ValidationError{Field: "category", Message: "is required"}
	case strings.TrimSpace(*category) == "":
		return out, &ValidationError{Field: "category", Message: "is required"}
	case brand == nil:
		return out, &ValidationError{Field: "brand", Message: "is required"}
	case strings.TrimSpace(*brand) == "":
		return out, &ValidationError{Field: "brand", Message: "is required"}
	case price == nil:
		return out, &ValidationError{Field: "price", Message: "is required"}
	case *price <= 0:
		return out, &ValidationError{Field: "price", Message: "must be greater than 0"}
	case stock == nil:
		return out, &ValidationError{Field: "stock", Message: "is required"}
	case *stock < 0:
		return out, &ValidationError{Field: "stock", Message: "must be 0 or greater"}
	default:
		out.Name = strings.TrimSpace(*name)
		out.SKU = strings.TrimSpace(*sku)
		out.Description = strings.TrimSpace(*description)
		out.Category = strings.TrimSpace(*category)
		out.Brand = strings.TrimSpace(*brand)
		out.Price = *price
		out.Stock = *stock
		if imageURL != nil {
			out.ImageURL = strings.TrimSpace(*imageURL)
		}
		out.IsActive = defaultIsActive
		if isActive != nil {
			out.IsActive = *isActive
		}
		return out, nil
	}
}
