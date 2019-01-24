package main

import (
	"encoding/json"
	"fmt"
	"search"

	"github.com/joaosoft/dbr"
)

type Person struct {
	IdPerson  int    `json:"id_person" db.read:"id_person"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age" db:"age"`
	Active    bool   `json:"active" db:"active"`
	IdAddress int    `json:"fk_address" db:"fk_address"`
}

type Address struct {
	IdAddress int    `json:"id_address" db:"id_address"`
	Street    string `json:"street" db:"street"`
	Number    int    `json:"number" db:"number"`
	Country   string `json:"country" db:"country"`
}

var db, _ = dbr.New()
var searcher, _ = search.New()

func main() {
	DeleteAll()

	Insert()

	Search()

	DeleteAll()
}

func Search() {
	data := []Person{}
	result, err := searcher.NewDatabaseSearch(db.Select("*").From("public.person")).Bind(&data).Page(2).Size(2).Exec()
	if err != nil {
		panic(err)
	}

	if result != nil {
		b, _ := json.MarshalIndent(result, "", "\t")
		fmt.Printf("\n\nResult: %s", string(b))
	}
}

func Insert() {
	fmt.Println("\n\n:: INSERT")

	address := Address{
		IdAddress: 1,
		Street:    "rua dos testes",
		Number:    1,
		Country:   "portugal",
	}
	if _, err := db.Insert().
		Into("public.address").
		Record(address).Exec(); err != nil {
		panic(err)
	}

	for i := 1; i <= 20; i++ {
		person := Person{
			IdPerson:  i,
			FirstName: "joao",
			LastName:  "ribeiro",
			Age:       i,
			IdAddress: 1,
		}

		if _, err := db.Insert().
			Into("public.person").
			Record(person).Exec(); err != nil {
			panic(err)
		}
	}
	fmt.Printf("\nINSERTED")
}

func DeleteAll() {
	fmt.Println("\n\n:: DELETE")

	if _, err := db.Delete().
		From("public.person").Exec(); err != nil {
		panic(err)
	}

	if _, err := db.Delete().
		From("public.address").Exec(); err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED")
}
