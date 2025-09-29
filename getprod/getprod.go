package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Структуры для DummyJSON
type DummyResponse struct {
	Products []DummyProduct `json:"products"`
}

type DummyProduct struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Images      []string `json:"images"`
}

func fetchProducts() ([]DummyProduct, error) {
	url := "https://dummyjson.com/products?limit=100"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Go-Client")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}

	var dummyResp DummyResponse
	if err := json.NewDecoder(resp.Body).Decode(&dummyResp); err != nil {
		return nil, err
	}

	return dummyResp.Products, nil
}

func insertProducts(db *pgxpool.Pool, products []DummyProduct) error {
	ctx := context.Background()
	inserted := 0
	for i, p := range products {
		if p.Title == "" || len(p.Images) == 0 {
			continue
		}

		price := p.Price
		amount := rand.Intn(50) + 1

		_, err := db.Exec(
			ctx,
			`INSERT INTO products(name, price, description, amount, image_url) 
			 VALUES ($1,$2,$3,$4,$5)`,
			p.Title,
			price,
			p.Description,
			amount,
			p.Images[0],
		)
		if err != nil {
			log.Printf("❌ Failed to insert product %d: %v\n", i+1, err)
			continue
		}

		inserted++
		fmt.Printf("✅ [%d] Added: %s | Price: $%.2f | Amount: %d\n", inserted, p.Title, price, amount)
	}

	fmt.Printf("\n🎉 Successfully inserted %d products into the database!\n", inserted)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dsn := "postgres://market:12345@localhost:5432/market?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()
	fmt.Println("Connected to the database")

	products, err := fetchProducts()
	if err != nil {
		log.Fatal("Failed to fetch products:", err)
	}

	err = insertProducts(pool, products)
	if err != nil {
		log.Fatal("Failed to insert products:", err)
	}
}
