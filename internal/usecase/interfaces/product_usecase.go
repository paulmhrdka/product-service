package interfaces

import (
	"context"
	"product-service/internal/domain"
	"product-service/internal/repository/elasticsearch"
)

type ProductUseCase interface {
	Search(ctx context.Context, params elasticsearch.SearchParams) (*domain.SearchProductResponse, error)
}
