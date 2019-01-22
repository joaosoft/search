package main

import (
	"fmt"

	"github.com/joaosoft/dbr"
)

type Person struct {
	IdPerson  int    `json:"id_person" db.read:"id_person"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Age       int    `json:"age" db:"age"`
	IdAddress *int   `json:"fk_address" db:"fk_address"`
}

type Address struct {
	IdAddress int    `json:"id_address" db:"id_address"`
	Street    string `json:"street" db:"street"`
	Number    int    `json:"number" db:"number"`
	Country   string `json:"country" db:"country"`
}

var db, _ = dbr.New()

func main() {
	DeleteAll()

	Insert()

	DeleteAll()
}

func Insert() {
	fmt.Println("\n\n:: INSERT")

	person := Person{
		FirstName: "joao",
		LastName:  "ribeiro",
		Age:       30,
	}

	stmt := db.Insert().
		Into(dbr.Field("public.person").As("new_name")).
		Record(person)

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nSAVED PERSON: %+v", person)
}

func DeleteAll() {
	fmt.Println("\n\n:: DELETE")

	stmt := db.Delete().
		From("public.person")

	query, err := stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	stmt = db.Delete().
		From("public.address")

	query, err = stmt.Build()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nQUERY: %s", query)

	_, err = stmt.Exec()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nDELETED ALL")
}
