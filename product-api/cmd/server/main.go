package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"

	model "products/internal"

	_ "github.com/mattn/go-sqlite3"

	"net/http"
)

type Service struct {
	database *sql.DB
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)

	dbPath := "./products.db"
	database, err := sql.Open("sqlite3",
		fmt.Sprintf("file:%s?mode=ro&"+
			"_journal=MEMORY&"+
			"_query_only=true", dbPath))
	if err != nil {
		panic(err)
	}
	service := &Service{database: database}

	http.Handle("/", middleware(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/random/", http.StatusTemporaryRedirect)
		})))
	http.Handle("/random/", middleware(http.HandlerFunc(service.serveProducts)))

	// Start server
	log.Info().Msg("Listening on port " + Port())
	log.Fatal().Err(http.ListenAndServe(":"+Port(), nil)).Msg("Can't start service")

}

func (s *Service) serveProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := s.database.QueryContext(r.Context(),
		"SELECT ProductID, Name, Type, "+
			"Price, Upc, Description, "+
			"Brand, Model, URL, "+
			"Image, Category FROM products "+
			"ORDER BY RANDOM() "+
			"LIMIT 10")

	if err != nil {
		sendErr(w, err)
		return
	}

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err = rows.Scan(&p.ProductID, &p.Name, &p.Type,
			&p.Price, &p.Upc, &p.Description,
			&p.Brand, &p.Model, &p.URL,
			&p.Image, &p.Category)
		if err != nil {
			sendErr(w, err)
			return
		}
		logger := log.Ctx(r.Context())
		logger.Info().Str("productID", p.ProductID).Str("brand", p.Brand).Msg("Product")
		products = append(products, p)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(products)
	if err != nil {
		sendErr(w, err)
		return
	}
}

func sendErr(w http.ResponseWriter, err error) {
	http.Error(w, "Error retrieving products", http.StatusInternalServerError)
	fmt.Println(err)
}
