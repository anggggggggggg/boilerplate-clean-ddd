package deposit

import "context"

// Repository adalah kontrak (interface) yang mendefinisikan operasi persistence untuk Deposit.
// Layer domain hanya tahu interface-nya, tidak tahu implementasinya (bisa PostgreSQL, MongoDB, dll).
type Repository interface {
	Create(ctx context.Context, deposit *Deposit) error
	FindByID(ctx context.Context, id string) (*Deposit, error)
	FindByUserID(ctx context.Context, userID string) ([]*Deposit, error)
	Update(ctx context.Context, deposit *Deposit) error
	Delete(ctx context.Context, id string) error
}
