package model

type (
	ProductRequest struct {
		Name  string  `json:"name" validate:"required"`
		Price float64 `json:"price" validate:"required"`
		Qty   int     `json:"qty" validate:"required"`
	}

	ProductResponse struct {
		Id        int     `json:"id"`
		Name      string  `json:"name"`
		Price     float64 `json:"price"`
		Qty       int     `json:"qty"`
		CreatedAt string  `json:"created_at,omitempty"`
	}
)
