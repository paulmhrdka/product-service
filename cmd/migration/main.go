package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/xuri/excelize/v2"
)

const (
	indexName     = "products"
	excelFilePath = "data/products.xlsx"
	sheetName     = "Sheet1"
	batchSize     = 1000
)

// Product mapping for Elasticsearch
const mapping = `
{
    "mappings": {
        "properties": {
            "id": { "type": "keyword" },
            "product_name": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
                "analyzer": "ngram_token_analyzer",
                "search_analyzer": "search_term_analyzer"
            },
            "drug_generic": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
                "analyzer": "ngram_token_analyzer",
                "search_analyzer": "search_term_analyzer"
            },
            "company": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
                "analyzer": "ngram_token_analyzer",
                "search_analyzer": "search_term_analyzer"
            },
            "created_at": { "type": "date" },
            "updated_at": { "type": "date" }
        }
    },
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 1,
        "analysis": {
            "analyzer": {
                "custom_analyzer": {
                    "type": "custom",
                    "tokenizer": "standard",
                    "filter": ["lowercase", "asciifolding"]
                }
            }
        }
    }
}
`

type Product struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	DrugGeneric string `json:"drug_generic"`
	Company     string `json:"company"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func main() {
	// Initialize Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Create or update index with mapping
	if err := createIndex(es); err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	// Read Excel file
	products, err := readExcel()
	if err != nil {
		log.Fatalf("Error reading Excel file: %s", err)
	}

	// Bulk import products
	if err := bulkImport(es, products); err != nil {
		log.Fatalf("Error importing products: %s", err)
	}

	log.Printf("Successfully imported %d products", len(products))
}

func createIndex(es *elasticsearch.Client) error {
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		res, err = es.Indices.Create(
			indexName,
			es.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("error creating index: %s", res.String())
		}
	}

	return nil
}

func readExcel() ([]Product, error) {
	f, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var products []Product
	// Skip header row
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 4 { // We need all columns
			log.Printf("Warning: Row %d has incomplete data, skipping", i+1)
			continue
		}

		// Clean and validate data
		id := strings.TrimSpace(row[0])
		productName := strings.TrimSpace(row[1])
		drugGeneric := strings.TrimSpace(row[2])
		company := strings.TrimSpace(row[3])

		// Skip if essential fields are empty
		if id == "" || productName == "" {
			log.Printf("Warning: Row %d has empty essential fields, skipping", i+1)
			continue
		}

		now := time.Now().Format(time.RFC3339)

		product := Product{
			ID:          id,
			ProductName: productName,
			DrugGeneric: drugGeneric,
			Company:     company,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		products = append(products, product)
	}

	return products, nil
}

func bulkImport(es *elasticsearch.Client, products []Product) error {
	var buf strings.Builder

	log.Printf("Starting bulk import of %d products...", len(products))
	successCount := 0

	for i, product := range products {
		// Create bulk operation metadata
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    product.ID,
			},
		}

		// Add metadata and document to buffer
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return err
		}
		if err := json.NewEncoder(&buf).Encode(product); err != nil {
			return err
		}

		// Execute bulk operation when batch size is reached or on last item
		if (i+1)%batchSize == 0 || i == len(products)-1 {
			res, err := es.Bulk(
				strings.NewReader(buf.String()),
				es.Bulk.WithContext(context.Background()),
				es.Bulk.WithIndex(indexName),
			)
			if err != nil {
				return err
			}

			// Parse response for detailed error reporting
			var bulkResponse map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
				return err
			}
			res.Body.Close()

			if res.IsError() {
				// Handle bulk errors
				if items, ok := bulkResponse["items"].([]interface{}); ok {
					for _, item := range items {
						if idx, ok := item.(map[string]interface{})["index"].(map[string]interface{}); ok {
							if errMsg, exists := idx["error"]; exists {
								log.Printf("Error importing product %s: %v", idx["_id"], errMsg)
							} else {
								successCount++
							}
						}
					}
				}
			} else {
				successCount += batchSize
			}

			// Clear buffer for next batch
			buf.Reset()

			// Log progress
			log.Printf("Processed %d/%d products...", i+1, len(products))
		}
	}

	log.Printf("Bulk import completed. Successfully imported %d products", successCount)
	return nil
}
