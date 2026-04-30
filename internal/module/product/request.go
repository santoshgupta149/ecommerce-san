package product

type createProductRequest struct {
	Name        *string  `json:"name"`
	SKU         *string  `json:"sku"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Brand       *string  `json:"brand"`
	Price       *float64 `json:"price"`
	Stock       *int64   `json:"stock"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

type updateProductRequest struct {
	Name        *string  `json:"name"`
	SKU         *string  `json:"sku"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Brand       *string  `json:"brand"`
	Price       *float64 `json:"price"`
	Stock       *int64   `json:"stock"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}
