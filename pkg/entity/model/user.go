package model

type (
	RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	LoginResponse struct {
		Token   string `json:"token"`
		Expired int64  `json:"expired"`
	}

	UserResponse struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
	}

	ValidateResponse struct {
		UserResponse
		Status bool
	}
)
