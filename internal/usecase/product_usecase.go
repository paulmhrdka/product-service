package usecase

import (
	"context"
	"math"
	"product-service/internal/domain"
	"product-service/internal/repository/elasticsearch"
)

type ProductUseCase struct {
	productRepo *elasticsearch.ProductRepository
}

func NewProductUseCase(productRepo *elasticsearch.ProductRepository) *ProductUseCase {
	return &ProductUseCase{productRepo: productRepo}
}

func (uc *ProductUseCase) Search(ctx context.Context, params elasticsearch.SearchParams) (*domain.SearchProductResponse, error) {
	// Get products and total count from repository
	products, totalItems, err := uc.productRepo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	// Return default value for empty result
	if products == nil {
		products = []domain.Product{}
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.Size)))

	pagination := domain.Pagination{
		CurrentPage: params.Page,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
		PageSize:    params.Size,
	}

	response := &domain.SearchProductResponse{
		Results:    products,
		Pagination: pagination,
	}

	return response, nil
}
