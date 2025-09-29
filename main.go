package main

import (
	"context"
	"log"
	"net/http"

	"github.com/akrikkk/market/product"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://market:12345@localhost:5432/market?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("Connected to the database")

	http.HandleFunc("/", product.ProductHandler(pool))
	http.Handle("/viewe/", http.StripPrefix("/viewe/", http.FileServer(http.Dir("viewe"))))

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
