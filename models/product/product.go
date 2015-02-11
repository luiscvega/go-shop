package product

import (
	"database/sql"

	"github.com/luiscvega/model"
)

var DB *sql.DB

type Product struct {
	Id    int    `column:"id"`
	Name  string `column:"name"`
	Price int    `column:"price"`
}

func Fetch(p *Product) error {
	err := model.Fetch("products", p, DB)
	if err != nil {
		return err
	}

	return nil
}

func All() ([]Product, error) {
	products := make([]Product, 0)

	err := model.All("products", &products, DB)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func Create(p *Product) error {
	id, err := model.Create("products", *p, DB)
	if err != nil {
		return err
	}

	p.Id = id

	return nil
}

func Update(p Product) error {
	err := model.Update("products", p, DB)
	if err != nil {
		return err
	}

	return nil
}

func Delete(id int) error {
	err := model.Delete("products", id, DB)
	if err != nil {
		return err
	}

	return nil
}
