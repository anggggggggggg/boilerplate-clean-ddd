package http

import (
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/dto"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/internal/usecase"
	"github.com/anggggggggggg/boilerplate-clean-and-ddd/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// DepositHandler adalah HTTP handler untuk domain deposit.
// Tanggung jawabnya hanya 3 hal:
//  1. Parse HTTP request menjadi DTO
//  2. Panggil use case
//  3. Tulis HTTP response
//
// Tidak ada business logic di sini — handler adalah "gerbang" ke use case.
type DepositHandler struct {
	depositUseCase usecase.DepositUseCase
}

func NewDepositHandler(depositUseCase usecase.DepositUseCase) *DepositHandler {
	return &DepositHandler{depositUseCase: depositUseCase}
}

// POST /api/deposits
func (h *DepositHandler) CreateDeposit(c *fiber.Ctx) error {
	var req dto.CreateDepositRequest
	if err := c.BodyParser(&req); err != nil {
		return response.WriteError(c, fiber.StatusBadRequest, "Bad Request", "body request tidak valid")
	}

	result, err := h.depositUseCase.CreateDeposit(c.Context(), &req)
	if err != nil {
		return response.WriteError(c, fiber.StatusUnprocessableEntity, "Gagal membuat deposit", err.Error())
	}

	return response.WriteSuccess(c, fiber.StatusCreated, "Deposit berhasil dibuat", result)
}

// GET /api/deposits/:id
func (h *DepositHandler) GetDeposit(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.WriteError(c, fiber.StatusBadRequest, "Bad Request", "id tidak boleh kosong")
	}

	result, err := h.depositUseCase.GetDeposit(c.Context(), id)
	if err != nil {
		return response.WriteError(c, fiber.StatusNotFound, "Deposit tidak ditemukan", err.Error())
	}

	return response.WriteSuccess(c, fiber.StatusOK, "Deposit berhasil diambil", result)
}

// GET /api/deposits/user/:user_id
func (h *DepositHandler) GetDepositsByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return response.WriteError(c, fiber.StatusBadRequest, "Bad Request", "user_id tidak boleh kosong")
	}

	results, err := h.depositUseCase.GetDepositsByUserID(c.Context(), userID)
	if err != nil {
		return response.WriteError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar deposit", err.Error())
	}

	return response.WriteSuccess(c, fiber.StatusOK, "Daftar deposit berhasil diambil", results)
}
