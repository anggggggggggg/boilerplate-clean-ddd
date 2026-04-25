package deposit

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// DepositStatus adalah value object yang merepresentasikan status deposit.
// Menggunakan tipe khusus agar compiler mencegah nilai status yang tidak valid.
type DepositStatus string

const (
	DepositStatusPending DepositStatus = "pending"
	DepositStatusSuccess DepositStatus = "success"
	DepositStatusFailed  DepositStatus = "failed"
)

// Deposit adalah aggregate root dari domain deposit.
// Semua business rule deposit ada di sini — tidak ada logic bisnis di layer lain.
type Deposit struct {
	ID        string
	UserID    string
	Amount    float64
	Status    DepositStatus
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewDeposit adalah factory function untuk membuat deposit baru.
// Setiap deposit baru selalu dimulai dengan status "pending" — ini adalah business rule.
func NewDeposit(userID string, amount float64, note string) (*Deposit, error) {
	if amount <= 0 {
		return nil, errors.New("amount harus lebih dari 0")
	}
	if userID == "" {
		return nil, errors.New("user_id tidak boleh kosong")
	}
	return &Deposit{
		ID:     uuid.New().String(),
		UserID: userID,
		Amount: amount,
		Status: DepositStatusPending,
		Note:   note,
	}, nil
}

// Approve adalah business method: menyetujui deposit.
// Hanya deposit dengan status pending yang dapat diapprove.
func (d *Deposit) Approve() error {
	if !d.IsPending() {
		return errors.New("hanya deposit berstatus pending yang dapat diapprove")
	}
	d.Status = DepositStatusSuccess
	return nil
}

// Reject adalah business method: menolak deposit.
func (d *Deposit) Reject() error {
	if !d.IsPending() {
		return errors.New("hanya deposit berstatus pending yang dapat ditolak")
	}
	d.Status = DepositStatusFailed
	return nil
}

// IsPending mengecek apakah deposit masih dalam status pending.
func (d *Deposit) IsPending() bool {
	return d.Status == DepositStatusPending
}
