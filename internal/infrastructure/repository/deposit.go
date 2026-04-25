package repository

import (
	"context"
	"time"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/domain/deposit"
	"gorm.io/gorm"
)

// depositModel adalah representasi tabel "deposits" di database.
// Sengaja dipisah dari domain entity agar domain tidak bergantung pada ORM atau tag GORM.
// Pattern ini disebut "anti-corruption layer" — melindungi domain dari detail infrastruktur.
type depositModel struct {
	ID        string    `gorm:"column:id;primaryKey;type:uuid"`
	UserID    string    `gorm:"column:user_id;not null"`
	Amount    float64   `gorm:"column:amount;not null"`
	Status    string    `gorm:"column:status;not null;default:'pending'"`
	Note      string    `gorm:"column:note"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (depositModel) TableName() string {
	return "deposits"
}

// depositRepository mengimplementasikan deposit.Repository menggunakan GORM.
// Ini satu-satunya file yang boleh tahu tentang GORM dan SQL.
type depositRepository struct {
	db *gorm.DB
}

// NewDepositRepository membuat instance depositRepository baru.
// Dikembalikan sebagai deposit.Repository (interface) sehingga caller tidak tahu concrete type-nya.
func NewDepositRepository(db *gorm.DB) deposit.Repository {
	return &depositRepository{db: db}
}

func (r *depositRepository) Create(ctx context.Context, d *deposit.Deposit) error {
	model := toModel(d)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	// Sync timestamps kembali ke domain entity setelah disimpan
	d.CreatedAt = model.CreatedAt
	d.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *depositRepository) FindByID(ctx context.Context, id string) (*deposit.Deposit, error) {
	var model depositModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return toDomain(&model), nil
}

func (r *depositRepository) FindByUserID(ctx context.Context, userID string) ([]*deposit.Deposit, error) {
	var models []depositModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	deposits := make([]*deposit.Deposit, 0, len(models))
	for i := range models {
		deposits = append(deposits, toDomain(&models[i]))
	}
	return deposits, nil
}

func (r *depositRepository) Update(ctx context.Context, d *deposit.Deposit) error {
	return r.db.WithContext(ctx).Save(toModel(d)).Error
}

func (r *depositRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&depositModel{}).Error
}

// toModel mengkonversi domain entity → database model (untuk write operations)
func toModel(d *deposit.Deposit) *depositModel {
	return &depositModel{
		ID:        d.ID,
		UserID:    d.UserID,
		Amount:    d.Amount,
		Status:    string(d.Status),
		Note:      d.Note,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

// toDomain mengkonversi database model → domain entity (untuk read operations)
func toDomain(m *depositModel) *deposit.Deposit {
	return &deposit.Deposit{
		ID:        m.ID,
		UserID:    m.UserID,
		Amount:    m.Amount,
		Status:    deposit.DepositStatus(m.Status),
		Note:      m.Note,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
