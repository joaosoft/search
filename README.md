search
================

[![Build Status](https://travis-ci.org/joaosoft/search.svg?branch=master)](https://travis-ci.org/joaosoft/search) | [![codecov](https://codecov.io/gh/joaosoft/search/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/search) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/search)](https://goreportcard.com/report/github.com/joaosoft/search) | [![GoDoc](https://godoc.org/github.com/joaosoft/search?status.svg)](https://godoc.org/github.com/joaosoft/search)

A simple tool that allows to you to search with pagination on database and elastic.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* database search
* elastic search (under development)

## Dependency Management
>### Dependency

Project dependencies are managed using Dependency. Read more about [Dependency](https://github.com/joaosoft/dependency).
* Get dependency manager: `go get github.com/joaosoft/dependency`

###### Commands
* Install dependencies: `dependency get`
* Update dependencies: `dependency update`
* Reset dependencies: `dependency reset`
* Add dependencies: `dependency add <dependency>`
* Remove dependencies: `dependency remove <dependency>`

>### Go
```
go get github.com/joaosoft/search
```

>### Configuration
>>#### master / slave
```
{
  "search": {
    "log": {
      "level": "info"
    }
  },
  "dbr": {
    "read_db": {
      "driver": "postgres",
      "datasource": "postgres://user:password@localhost:7000/postgres?sslmode=disable&search_path=public"
    },
    "write_db": {
      "driver": "postgres",
      "datasource": "postgres://user:password@localhost:7100/postgres?sslmode=disable&search_path=public"
    },
    "log": {
      "level": "info"
    }
  }
}
```

## Usage 
This examples are available in the project at [search/examples](https://github.com/joaosoft/search/tree/master/examples)

```go
type Person struct {
	IdPerson  int    `json:"id_person" db:"id_person"`
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

	result, err := searcher.NewDatabaseSearch(
		db.Select("*").
			From("public.person").
			OrderAsc("id_person")).
		Bind(&[]Person{}).
		Path("http://teste.pt").
		Page(1).
		Size(3).
		Metadata("my-meta",
			db.Select("*").
				From("public.person").
				OrderAsc("id_person"),
			&[]Person{}).
		Exec()

	if err != nil {
		panic(err)
	}

	if result != nil {
		b, _ := json.MarshalIndent(result, "", "\t")
		fmt.Printf("\n\nSearch: %s", string(b))
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
```

> ##### Result:
```
:: DELETE

DELETED

:: INSERT

INSERTED

Search: {
	"result": [
		{
			"id_person": 1,
			"first_name": "joao",
			"last_name": "ribeiro",
			"age": 1,
			"active": false,
			"fk_address": 1
		},
		{
			"id_person": 2,
			"first_name": "joao",
			"last_name": "ribeiro",
			"age": 2,
			"active": false,
			"fk_address": 1
		},
		{
			"id_person": 3,
			"first_name": "joao",
			"last_name": "ribeiro",
			"age": 3,
			"active": false,
			"fk_address": 1
		}
	],
	"metadata": {
		"my-meta": [
			{
				"id_person": 1,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 1,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 2,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 2,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 3,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 3,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 4,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 4,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 5,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 5,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 6,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 6,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 7,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 7,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 8,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 8,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 9,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 9,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 10,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 10,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 11,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 11,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 12,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 12,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 13,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 13,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 14,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 14,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 15,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 15,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 16,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 16,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 17,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 17,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 18,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 18,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 19,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 19,
				"active": false,
				"fk_address": 1
			},
			{
				"id_person": 20,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 20,
				"active": false,
				"fk_address": 1
			}
		]
	},
	"pagination": {
		"first": null,
		"previous": null,
		"next": "http://teste.pt?page=2\u0026size=3",
		"last": "http://teste.pt?page=7\u0026size=3"
	}
}
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
