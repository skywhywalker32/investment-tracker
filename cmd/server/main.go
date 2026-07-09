package main

import (
	"database/sql"
	"log"

	stockhttp "investment-tracker/internal/stock/entrypoint/http"
	stockservice "investment-tracker/internal/stock/service"
	stockstore "investment-tracker/internal/stock/store"

	txhttp "investment-tracker/internal/transaction/entrypoint/http"
	txservice "investment-tracker/internal/transaction/service"
	txstore "investment-tracker/internal/transaction/store"

	portfoliohttp "investment-tracker/internal/portfolio/entrypoint/http"
	portfolioservice "investment-tracker/internal/portfolio/service"
	portfoliostore "investment-tracker/internal/portfolio/store"

	"github.com/gofiber/fiber/v3"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5433/stock_service?sslmode=disable"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Не удалось инициализировать пул баз данных: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL!")

	app := fiber.New()

	skStore := stockstore.New(db)
	txStore := txstore.New(db)
	prStore := portfoliostore.New(db)

	skSvc := stockservice.New(skStore)
	txSvc := txservice.New(txStore)
	prSvc := portfolioservice.New(prStore)

	skHandler := stockhttp.NewHandler(skSvc)
	txHandler := txhttp.NewHandler(txSvc)
	prHandler := portfoliohttp.NewHandler(prSvc)

	app.Get("/api/v1/stocks", skHandler.GetTickers)
	app.Get("/api/v1/stocks/:ticker", skHandler.GetByTicker)
	app.Post("/api/v1/stocks", skHandler.CreateStock)
	app.Get("/api/v1/stocks/:ticker/prices", skHandler.GetPriceHistory)

	app.Post("/api/v1/transactions", txHandler.PostTransaction)

	app.Get("/api/v1/users/:id/portfolio", prHandler.GetPortfolio)

	log.Fatal(app.Listen(":3030"))
}
