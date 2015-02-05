package product

import (
	"database/sql"
	"github.com/luiscvega/model"
	"log"
)

var DB *sql.DB

type Product struct {
	Id    int    `column:"id"`
	Name  string `column:"name"`
	Price int    `column:"price"`
}

func All() ([]Product, error) {
	products := make([]Product, 0)

	rows, err := DB.Query("SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
                p := Product{}
		err := rows.Scan(&p.Id, &p.Name, &p.Price)
		if err != nil {
			log.Fatal(err)
		}
                products = append(products, p)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
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
