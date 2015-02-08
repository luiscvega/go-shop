package main

import (
	"database/sql"
	"fmt"
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

	mux.Put("/products/:id", func(c *cuba.Context) {
		var p product.Product
		p.Id, _ = strconv.Atoi(c.Params["id"])
		p.Name = c.R.FormValue("name")
		p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
		product.Update(&p)

		c.Redirect("/")
	})

	mux.Delete("/products/:id", func(c *cuba.Context) {
		product.Delete(c.Params["id"])

		c.Redirect("/")
	})

	mux.Post("/products", func(c *cuba.Context) {
		var p product.Product
		p.Name = c.R.FormValue("name")
		p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
		product.Create(&p)

		c.Redirect("/")
	})

	mux.Get("/about/:first_name/boom/:last_name", func(c *cuba.Context) {
		c.Render("about", map[string]string{
			"FirstName": c.Params["first_name"],
			"LastName":  c.Params["last_name"]})
	})

	mux.Get("/", func(c *cuba.Context) {
		products, err := product.All()
		if err != nil {
			panic(err)
		}

		c.Render("index", products)
	})

	for method, routes := range mux.Table() {
		fmt.Println("METHOD:", method)

		for _, route := range routes {
			fmt.Println(route)
		}

		fmt.Println("=====================================================================================")
	}

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", mux)
}
