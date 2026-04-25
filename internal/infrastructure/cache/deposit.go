package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/dto"
	"github.com/redis/go-redis/v9"
)

// depositTTL adalah durasi data deposit disimpan di cache.
// Setelah TTL habis, Redis otomatis menghapus key tersebut.
const depositTTL = 5 * time.Minute

// depositCache mengimplementasikan usecase.DepositCache menggunakan Redis.
// Struct ini tidak mengimport package usecase \u2014 Go structural typing yang memastikan
// implementasi ini memenuhi interface tanpa perlu deklarasi eksplisit.
type depositCache struct {
	client *redis.Client
}

// NewDepositCache membuat instance depositCache baru.
// Dikembalikan sebagai concrete type agar factory bisa assign ke interface.
func NewDepositCache(client *redis.Client) *depositCache {
	return &depositCache{client: client}
}

// Get mengambil deposit dari cache berdasarkan ID.
// Mengembalikan error jika key tidak ditemukan (redis.Nil) atau data rusak.
// Caller (use case) yang menentukan apakah error ini kritikal atau tidak.
func (c *depositCache) Get(ctx context.Context, id string) (*dto.DepositResponse, error) {
	data, err := c.client.Get(ctx, depositKey(id)).Bytes()
	if err != nil {
		return nil, err // termasuk redis.Nil jika key tidak ada
	}

	var res dto.DepositResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("gagal unmarshal deposit dari cache: %w", err)
	}

	return &res, nil
}

// Set menyimpan deposit ke cache dengan TTL yang sudah ditentukan.
func (c *depositCache) Set(ctx context.Context, d *dto.DepositResponse) error {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("gagal marshal deposit untuk cache: %w", err)
	}

	return c.client.Set(ctx, depositKey(d.ID), data, depositTTL).Err()
}

// Delete menghapus deposit dari cache. Dipanggil ketika data diupdate/dihapus.
func (c *depositCache) Delete(ctx context.Context, id string) error {
	return c.client.Del(ctx, depositKey(id)).Err()
}

// depositKey menghasilkan Redis key yang konsisten untuk deposit.
// Format: "deposit:{id}"
func depositKey(id string) string {
	return fmt.Sprintf("deposit:%s", id)
}
