package model

type (
	OrderRequest struct {
		UserId    int64   `json:"user_id" validate:"required"`
		ProductId int64   `json:"product_id" validate:"required"`
		Amount    float64 `json:"amount" validate:"required"`
		Quantity  int64   `json:"quantity" validate:"required"`
	}

	OrderResponse struct {
		OrderId     int64   `json:"order_id"`
		UserId      int64   `json:"user_id"`
		UserName    string  `json:"user_name"`
		ProductId   int64   `json:"product_id"`
		ProductName string  `json:"product_name"`
		Amount      float64 `json:"amount" `
		Status      string  `json:"status"`
	}

	PaymentRequest struct {
		UserId    int64   `json:"user_id" validate:"required"`
		ProductId int64   `json:"product_id" validate:"required"`
		Amount    float64 `json:"amount" validate:"required"`
	}
)
