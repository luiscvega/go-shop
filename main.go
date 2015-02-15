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

func main() {
	var err error
	product.DB, err = sql.Open("postgres", "dbname=shop")
	if err != nil {
		panic(err)
	}
	defer product.DB.Close()

	mux := cuba.New()

	mux.On("products", func(mux *cuba.Mux) {

		mux.Post("", func(c *cuba.Context) error {
			var p product.Product
			p.Name = c.R.FormValue("name")
			p.Price, _ = strconv.Atoi(c.R.FormValue("price"))

			err := product.Create(&p)
			if err != nil {
				return err
			}

			c.Redirect(fmt.Sprintf("/products/%d", p.Id))

			return nil
		})

		mux.Get(":id", func(c *cuba.Context) error {
			var p product.Product
			p.Id, _ = strconv.Atoi(c.Params["id"])

			err := product.Fetch(&p)
			if err != nil {
				return err
			}

			c.Render("products/show", p)

			return nil
		})

		mux.Put(":id", func(c *cuba.Context) error {
			var p product.Product
			p.Id, _ = strconv.Atoi(c.Params["id"])
			p.Name = c.R.FormValue("name")
			p.Price, _ = strconv.Atoi(c.R.FormValue("price"))

			err := product.Update(p)
			if err != nil {
				return err
			}

			c.Redirect("/")

			return nil
		})

		mux.Delete(":id", func(c *cuba.Context) error {
			id, _ := strconv.Atoi(c.Params["id"])
			err := product.Delete(id)

			if err != nil {
				return err
			}

			c.Redirect("/")

			return nil
		})

		mux.Get("/", func(c *cuba.Context) error {
			products, err := product.All()
			if err != nil {
				return err
			}

			c.Render("products/index", products)

			return nil
		})

	})

	mux.Get("/", func(c *cuba.Context) error {
		c.Render("index", nil)

		return nil
	})

	for _, routes := range mux.Routes() {
		fmt.Println(routes)
	}

	c := chain.New(mux)
	c.Use(ServeStaticFiles)
	c.Use(MethodOverride)

	fmt.Println("Starting...")
	http.ListenAndServe(":8080", c)
}

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

func MethodOverride(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			if r.FormValue("_method") == "PUT" {
				r.Method = "PUT"
			}

			if r.FormValue("_method") == "DELETE" {
				r.Method = "DELETE"
			}
		}

		h.ServeHTTP(w, r)
	})
}
