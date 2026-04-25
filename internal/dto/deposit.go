package dto

// CreateDepositRequest adalah struktur data yang diterima dari client (HTTP request body)
type CreateDepositRequest struct {
	UserID string  `json:"user_id" validate:"required"`
	Amount float64 `json:"amount"  validate:"required,gt=0"`
	Note   string  `json:"note"`
}

// DepositResponse adalah struktur data yang dikirim ke client (HTTP response body)
type DepositResponse struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Note      string  `json:"note"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
