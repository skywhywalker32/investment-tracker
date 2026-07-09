package http

import (
	"investment-tracker/internal/portfolio/service"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetPortfolio(c fiber.Ctx) error {
	idStr := c.Params("id")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат ID пользователя")
	}

	portfolio, err := h.svc.GetPortfolio(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить портфель пользователя")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":   userID,
		"portfolio": portfolio,
	})
}
