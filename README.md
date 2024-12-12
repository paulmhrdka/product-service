# Product Service API

A RESTful API service for product search using Elasticsearch. This service provides powerful search capabilities with features like fuzzy matching, pagination, and sorting.

## Features

- Product search with fuzzy matching
- Pagination support
- Sorting by relevance and product name
- Elasticsearch integration
- Swagger documentation
- Clean architecture implementation

## Prerequisites

Before running this project, make sure you have the following installed:
- Go 1.23 or later [(Download here)](https://go.dev/dl/)
- Elasticsearch 8.16.0 [(Download here)](https://www.elastic.co/downloads/elasticsearch)
- `swag` command line tool for Swagger generation [(Official Documentation)](https://github.com/swaggo/swag)
- Make (for running Makefile commands)

## Installation

1. Install the required Go package for Swagger:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Clone the repository:
```bash
git clone https://github.com/paulmhrdka/product-service.git
cd product-service
```

3. Install dependencies:
```bash
go mod download
```

4. Start Elasticsearch locally on port 9200

5. Create .env file in project root:
```bash
SERVER_PORT=8080
ES_URLS=http://localhost:9200
```

## Usage

The project includes several make commands for common operations:

1. Generate Swagger documentation:
```bash
make swagger
```

2. Run the application:
```bash
make run
```

3. Run data migration (to import products from Excel):
```bash
make migrate
```

4. Build the application:
```bash
make build
```

## API Endpoints

### Search Products
```
GET /api/v1/products/search
```

Query Parameters:
- `q` (optional): Search query string
- `page` (optional): Page number (default: 1)
- `size` (optional): Items per page (default: 10, max: 100)

Example Request:
```bash
curl "http://localhost:8080/api/v1/products/search?q=paracetamol&page=1&size=10"
```

Example Response:
```json
{
    "success": true,
    "message": "Products retrieved successfully",
    "data": [
        {
            "id": "1",
            "product_name": "Paracetamol 500mg",
            "drug_generic": "Paracetamol",
            "company": "Company Name",
            "score": 0.9
        }
    ],
}
```

## Project Structure

```
.
├── api/
│   └── swagger/           # Swagger documentation
├── cmd/
│   ├── api/              # Application entrypoint
│   └── migration/        # Data migration tool
├── config/               # Configuration
├── internal/
│   ├── domain/          # Domain models
│   ├── delivery/        # HTTP handlers and routing
│   ├── repository/      # Data access layer
│   └── usecase/         # Business logic
├── pkg/                 # Shared packages
├── Makefile
└── README.md
```

## Data Migration

To import product data from Excel:

1. Prepare your Excel file with the following columns:
   - id
   - product_name
   - drug_generic
   - company

2. Place the Excel file in the `data` directory as `products.xlsx`

3. Run the migration:
```bash
make migrate
```

## API Documentation

After starting the application, you can access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Development

To add new features or modify existing ones:

1. Create or modify handlers in `internal/delivery/handler/`
2. Update business logic in `internal/usecase/`
3. Add new repository methods in `internal/repository/` if needed
4. Update Swagger annotations and regenerate documentation
5. Add tests for new functionality

## Error Handling

The API uses standardized error responses:

- 200: Successful operation
- 400: Bad request
- 422: Validation error
- 500: Internal server error

Example error response:
```json
{
    "success": false,
    "message": "Validation failed",
    "error": [
        {
            "field": "size",
            "message": "page size cannot exceed 100"
        }
    ],
}
```