package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"products-frontend/lib/config"
	"products-frontend/lib/run"
	"products-frontend/product"

	"github.com/Masterminds/sprig"
)

type App struct {
	product *product.Client
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	client, err := product.NewClient(os.Getenv("PRODUCTS_URL"))
	if err != nil {
		panic(err)
	}

	app := &App{product: client}
	http.HandleFunc("/", app.serveIndex)

	// Start server
	log.Println("Listening on port " + Port())
	log.Fatal(http.ListenAndServe(":"+Port(), nil))
}

// serveIndex returns the index.html file
func (app *App) serveIndex(
	w http.ResponseWriter, r *http.Request) {

	type IndexPage struct {
		Products *[]product.Product
	}

	ps, err := app.product.GetProducts(r.Context())
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		http.Error(w, "Error retrieving products",
			http.StatusInternalServerError)
		return
	}

	// Render page template
	tpl := template.Must(
		template.New("index.html").
			Funcs(sprig.FuncMap()).
			ParseFiles("web/index.html"))
	tpl.Execute(w, &IndexPage{ps})
}

// serveIndex returns the index.html file
func (app *App) deleteSelf(
	w http.ResponseWriter, r *http.Request) {
	client, _ := run.NewClient()

	service := os.Getenv("K_SERVICE")

	client.DeleteSelf(config.Region(), service)
}
