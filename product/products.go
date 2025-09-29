package product

import (
	"context"
	"log"
	"net/http"

	"html/template"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID          int
	Name        string
	Price       float64
	Description string
	Amount      int
	Image_url   string
}

func GetAllProducts(db *pgxpool.Pool) ([]Product, error) {
	rows, err := db.Query(context.Background(), "SELECT id, name, price, description, amount, image_url FROM products")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Amount, &p.Image_url)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func ProductHandler(db *pgxpool.Pool) http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("./template/products.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := GetAllProducts(db)
		if err != nil {
			http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, products)
	}
}
