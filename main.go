package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"

	"./chain"
	"./cuba"
	"./models/product"
)

func ServeStaticFiles(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matched, _ := regexp.MatchString("^/public/(css)", r.URL.Path)
		if matched {
			http.StripPrefix("/public/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
			return
		}

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

	mux.Put("/products/:id", func(c *cuba.Context) error {
		var p product.Product
		p.Id, _ = strconv.Atoi(c.Params["id"])
		p.Name = c.R.FormValue("name")
		p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
		product.Update(&p)

		c.Redirect("/")

		return nil
	})

	mux.Delete("/products/:id", func(c *cuba.Context) error {
		product.Delete(c.Params["id"])

		c.Redirect("/")

		return nil
	})

	mux.Post("/products", func(c *cuba.Context) error {
		var p product.Product
		p.Name = c.R.FormValue("name")
		p.Price, _ = strconv.Atoi(c.R.FormValue("price"))
		product.Create(&p)

		c.Redirect("/")

		return nil
	})

	mux.Get("/about/:first_name/boom/:last_name", func(c *cuba.Context) error {
		c.Render("about", map[string]string{
			"FirstName": c.Params["first_name"],
			"LastName":  c.Params["last_name"]})

		return nil
	})

	mux.Get("/", func(c *cuba.Context) error {
		products, err := product.All()
		if err != nil {
			return err
		}

		c.Render("index", products)

		return nil
	})

	c := chain.New(mux)
	c.Use(ServeStaticFiles)

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", c)
}
