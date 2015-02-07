package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"./cuba"
	"./models/product"
)

func main() {
	db, err := sql.Open("postgres", "dbname=shop")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product.DB = db

	mux := cuba.New()

	mux.Add("/products/:id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			product.Delete(r.FormValue("params"))
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})

	mux.Add("/products", productsHandler)

	mux.Add("/", rootHandler)

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", mux)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var p product.Product
		p.Name = r.FormValue("name")
		p.Price, _ = strconv.Atoi(r.FormValue("price"))
		product.Create(&p)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	products, err := product.All()
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))
	tmpl.ExecuteTemplate(w, "index.html", products)
}
