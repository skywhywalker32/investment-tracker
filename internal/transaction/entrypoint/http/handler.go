package http

import (
	"errors"
	"investment-tracker/internal/models"
	"investment-tracker/internal/transaction/service"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) PostTransaction(c fiber.Ctx) error {
	var tx models.Transaction

	if err := c.Bind().Body(&tx); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат JSON-данных")
	}

	createdTx, err := h.svc.ExecuteTransaction(c.Context(), &tx)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientStocks) {
			return fiber.NewError(fiber.StatusBadRequest, "Недостаточно акций в портфеле для совершения продажи")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка при проведении сделки")
	}

	return c.Status(fiber.StatusCreated).JSON(createdTx)
}
