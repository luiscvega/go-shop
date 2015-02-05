package product

import (
	"database/sql"
	"os/exec"
	"testing"

	_ "github.com/lib/pq"
)

func setup() {
	exec.Command("createdb", "products-test").Run()
	exec.Command("psql", "-f", "products.sql", "products-test").Run()

	var err error
	DB, err = sql.Open("postgres", "dbname=products-test")
	if err != nil {
		panic(err)
	}
}

func teardown() {
	DB.Close()
	exec.Command("dropdb", "products-test").Run()
}

func TestCreate(t *testing.T) {
	setup()
	defer teardown()

	p := Product{Name: "luis's stuff", Price: 300}

	err := Create(&p)
	if err != nil {
		t.Error(err)
	}
	if p.Id != 1 {
		t.Error("id != 1")
	}
}
