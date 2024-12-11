package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"product-service/internal/domain"
	"product-service/pkg/utils"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductRepository struct {
	es *elasticsearch.Client
}

func NewProductRepository(es *elasticsearch.Client) *ProductRepository {
	return &ProductRepository{es: es}
}

type SearchParams struct {
	Query string
	Page  int
	Size  int
}

func (r *ProductRepository) Search(ctx context.Context, params SearchParams) ([]domain.Product, int64, error) {
	from := (params.Page - 1) * params.Size

	// default elasticsearch query
	baseQuery := map[string]interface{}{
		"match_all": map[string]interface{}{},
	}

	// handle if params query not empty
	if params.Query != "" {
		baseQuery = map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"product_name": map[string]interface{}{
								"query":     params.Query,
								"operator":  "and",
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match": map[string]interface{}{
							"drug_generic": map[string]interface{}{
								"query":     params.Query,
								"operator":  "and",
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match": map[string]interface{}{
							"company": map[string]interface{}{
								"query":     params.Query,
								"operator":  "and",
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"multi_match": map[string]interface{}{
							"query": params.Query,
							"fields": []string{
								"product_name^3",
								"drug_generic^2",
								"company",
							},
							"type": "phrase",
						},
					},
				},
				"minimum_should_match": 1,
			},
		}
	}

	searchQuery := map[string]interface{}{
		"query":            baseQuery,
		"from":             from,
		"size":             params.Size,
		"track_total_hits": true,
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"product_name.keyword": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}

	queryJson, _ := json.MarshalIndent(searchQuery, "", "  ")
	log.Printf("Search Query: %s", string(queryJson))

	searchBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("error marshaling query: %w", err)
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("products"),
		r.es.Search.WithBody(strings.NewReader(string(searchBody))),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("error performing search: %w", err)
	}
	defer res.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("error decoding response: %w", err)
	}

	// Get total hits
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid response format")
	}

	total := int64(0)
	if totalHits, ok := hits["total"].(map[string]interface{}); ok {
		if value, ok := totalHits["value"].(float64); ok {
			total = int64(value)
		}
	}

	// Parse hits
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits format")
	}

	products := make([]domain.Product, 0)
	for _, hit := range hitsArray {
		h, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := h["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		score, _ := h["_score"].(float64)

		product := domain.Product{
			ID:          utils.GetString(source, "id"),
			ProductName: utils.GetString(source, "product_name"),
			DrugGeneric: utils.GetString(source, "drug_generic"),
			Company:     utils.GetString(source, "company"),
			Score:       score,
		}
		products = append(products, product)
	}

	return products, total, nil
}
