package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"

	"./models/product"
)

func main() {
	db, err := sql.Open("postgres", "dbname=shop")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	product.DB = db

	http.HandleFunc("/products/new", func(w http.ResponseWriter, r *http.Request) {
		p := product.Product{}
		p.Name = r.FormValue("name")
		p.Price, _ = strconv.Atoi(r.FormValue("price"))
		product.Create(&p)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		products, err := product.All()
		if err != nil {
			panic(err)
		}

		tmpl := template.Must(template.ParseFiles("views/index.html"))
		tmpl.ExecuteTemplate(w, "index.html", products)
	})

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", nil)
}
