package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	"./cuba"
	"./models/product"
)

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("PATH:", r.URL.Path)

		h.ServeHTTP(w, r)
	})
}

func TimeLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("TIME:", time.Now())

		h.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	product.DB, err = sql.Open("postgres", "dbname=shop")
	if err != nil {
		panic(err)
	}
	defer product.DB.Close()

	mux := cuba.New()

	mux.Use(RequestLogger)
	mux.Use(TimeLogger)

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

	//for method, routes := range mux.Table() {
	//fmt.Println("METHOD:", method)

	//for _, route := range routes {
	//fmt.Println(route)
	//}

	//fmt.Println("=====================================================================================")
	//}

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", mux)
}
