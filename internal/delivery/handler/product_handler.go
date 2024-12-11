package handler

import (
	"product-service/internal/repository/elasticsearch"
	"product-service/internal/usecase"
	"product-service/pkg/response"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// @Summary Search products
// @Description Search products with pagination
// @Tags Products
// @Accept json
// @Produce json
// @Param q query string false "Search query (optional)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/products/search [get]
func (h *ProductHandler) Search(c echo.Context) error {
	// Get query parameters
	query := strings.TrimSpace(c.QueryParam("q"))

	// Parse pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil || size < 1 {
		size = 10
	}

	// Maximum page size to prevent overload
	if size > 100 {
		return response.ValidationError(c, []response.ErrorDetail{
			{
				Field:   "size",
				Message: "page size cannot exceed 100",
			},
		})
	}

	// Prepare search parameters
	params := elasticsearch.SearchParams{
		Query: query,
		Page:  page,
		Size:  size,
	}

	// Perform search
	result, err := h.productUseCase.Search(c.Request().Context(), params)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	// Handle empty results
	if len(result.Results) == 0 {
		return response.Success(c, []interface{}{}, "No products found")
	}

	// Return successful response with pagination
	return response.SuccessWithPagination(c,
		result.Results,
		response.PaginationMeta{
			CurrentPage: result.Pagination.CurrentPage,
			TotalPages:  result.Pagination.TotalPages,
			TotalItems:  result.Pagination.TotalItems,
			PageSize:    result.Pagination.PageSize,
		},
		"Products retrieved successfully",
	)
}
