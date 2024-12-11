package domain

type Product struct {
	ID          string  `json:"id"`
	ProductName string  `json:"product_name"`
	DrugGeneric string  `json:"drug_generic"`
	Company     string  `json:"company"`
	Score       float64 `json:"score"`
}

type SearchProductResponse struct {
	Results    []Product  `json:"results"`
	Pagination Pagination `json:"pagination"`
}
