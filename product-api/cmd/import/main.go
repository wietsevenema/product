package main

import (
	"compress/gzip"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/fnv"
	model "products/internal"

	_ "github.com/mattn/go-sqlite3"

	"log"

	"os"
	"strconv"
)

func main() {

	sourcePath := "./assets/products/products.json.gz"
	dbPath := "./products.db"

	_, err := os.Stat(dbPath)
	if err == nil {
		log.Printf("Cowardly refusing to overwrite %s", dbPath)
		os.Exit(1)
	}

	database, _ := sql.Open("sqlite3",
		dbPath)

	_, err = database.Exec("CREATE TABLE IF NOT EXISTS products " +
		"(" +
		"ProductID INTEGER PRIMARY KEY," +
		"Name TEXT," +
		"Type TEXT," +
		"Price INTEGER," +
		"Upc TEXT," +
		"Description TEXT," +
		"Brand TEXT," +
		"Model TEXT," +
		"URL TEXT," +
		"Image TEXT," +
		"Category TEXT" +
		")")
	panicOnError(err)

	file, err := os.Open(sourcePath)
	panicOnError(err)
	defer file.Close()

	r, err := gzip.NewReader(file)
	panicOnError(err)
	defer r.Close()

	decoder := json.NewDecoder(r)

	for decoder.More() {
		t, err := decoder.Token()
		panicOnError(err)
		if d, ok := t.(json.Delim); !ok || d != '[' {
			log.Fatalf("expected array start token, got: %s", d)
		}

		ctx := context.Background()

		log.Printf("Importing products from %s", sourcePath)
		i := 0
		for ; decoder.More(); i++ {
			if i > 5000 {
				break
			}
			var raw ProductRaw
			err = decoder.Decode(&raw)
			panicOnError(err)

			p := toProduct(raw)
			statement, err := database.Prepare(
				"INSERT INTO products ( " +
					"ProductID, Name, Type, " +
					"Price, Upc, Description, " +
					"Brand, Model, URL, " +
					"Image, Category) " +
					"VALUES (:ProductID, :Name, :Type, " +
					":Price, :Upc, :Description, " +
					":Brand, :Model, :URL, " +
					":Image, :Category)")
			panicOnError(err)

			_, err = statement.ExecContext(ctx,
				sql.Named("ProductID", p.ProductID),
				sql.Named("Name", p.Name),
				sql.Named("Type", p.Type),
				sql.Named("Price", p.Price),
				sql.Named("Upc", p.Upc),
				sql.Named("Description", p.Description),
				sql.Named("Brand", p.Brand),
				sql.Named("Model", p.Model),
				sql.Named("URL", p.URL),
				sql.Named("Image", p.Image),
				sql.Named("Category", p.Category))
			panicOnError(err)
		}
		log.Printf("Imported %d products", i)
		break
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductRaw struct {
	Sku         int        `json:"sku"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Price       float64    `json:"price"`
	Upc         string     `json:"upc"`
	Description string     `json:"description"`
	Brand       string     `json:"manufacturer"`
	Model       string     `json:"model"`
	URL         string     `json:"url"`
	Categories  []Category `json:"category"`
	Image       string     `json:"image"`
}

func toProduct(pr ProductRaw) *model.Product {
	return &model.Product{
		ProductID:   hash(pr.Sku),
		Name:        pr.Name,
		Type:        pr.Type,
		Price:       int(pr.Price * 100),
		Upc:         pr.Upc,
		Description: pr.Description,
		Brand:       pr.Brand,
		Model:       pr.Model,
		URL:         pr.URL,
		Image:       pr.Image,
		Category:    extractName(pr.Categories),
	}
}

func hash(inp int) string {
	h := fnv.New32a()
	h.Write([]byte(strconv.Itoa(inp)))
	return fmt.Sprintf("%d", h.Sum32())

}

func extractName(cs []Category) string {
	for _, c := range cs {
		return c.Name
	}
	return ""
}
