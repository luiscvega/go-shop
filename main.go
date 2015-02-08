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

	mux.Add("/products/:id", func(c *cuba.Context) {
		if c.R.Method == "PUT" {
			var p product.Product
			p.Id, _ = strconv.Atoi(c.Params["id"])
			p.Name = c.R.FormValue("name")
			p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
			product.Update(&p)
		}

		if c.R.Method == "DELETE" {
			product.Delete(c.Params["id"])
		}

		http.Redirect(c.W, c.R, "/", http.StatusFound)
	})

	mux.Add("/products", productsHandler)

	mux.Add("/", rootHandler)

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", mux)
}

func productsHandler(c *cuba.Context) {
	if c.R.Method == "POST" {
		var p product.Product
		p.Name = c.R.FormValue("name")
		p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
		product.Create(&p)
	}

	http.Redirect(c.W, c.R, "/", http.StatusFound)
}

func rootHandler(c *cuba.Context) {
	products, err := product.All()
	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles("views/index.html"))
	tmpl.ExecuteTemplate(c.W, "index.html", products)
}
