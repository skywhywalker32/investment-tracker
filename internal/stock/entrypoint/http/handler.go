package http

import (
	"investment-tracker/internal/models"
	"investment-tracker/internal/stock/service"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetTickers(c fiber.Ctx) error {
	tickers, err := h.svc.GetTickers()
	if err != nil {
		log.Printf("error fetching tickers from db: %v", err)

		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить список акций")
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{
			"tickers": tickers,
		})
}

func (h *Handler) GetByTicker(c fiber.Ctx) error {
	ticker := c.Params("ticker")

	stock, err := h.svc.GetByTicker(ticker)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить акцию по её имени")
	}

	return c.Status(fiber.StatusOK).JSON(stock)
}

func (h *Handler) CreateStock(c fiber.Ctx) error {
	var stock models.Stock

	if err := c.Bind().Body(&stock); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат данных")
	}

	if err := h.svc.CreateStock(&stock); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось сохранить акцию")
	}

	return c.Status(fiber.StatusCreated).JSON(stock)
}

func (h *Handler) GetPriceHistory(c fiber.Ctx) error {
	ticker := c.Params("ticker")
	periodStr := c.Query("period", "168h")

	period, err := time.ParseDuration(periodStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат периода (используйте например '24h', '168h')")
	}

	history, err := h.svc.GetPriceHistory(ticker, period)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить историю цен")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ticker":  ticker,
		"period":  period,
		"history": history,
	})
}
