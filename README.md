search
================

[![Build Status](https://travis-ci.org/joaosoft/search.svg?branch=master)](https://travis-ci.org/joaosoft/search) | [![codecov](https://codecov.io/gh/joaosoft/search/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/search) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/search)](https://goreportcard.com/report/github.com/joaosoft/search) | [![GoDoc](https://godoc.org/github.com/joaosoft/search?status.svg)](https://godoc.org/github.com/joaosoft/search)

A simple tool that allows you to search with pagination on database and elastic.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* database search
* elastic search

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
```
{
  "search": {
    "log": {
      "level": "error"
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
      "level": "error"
    }
  },
  "elastic": {
    "endpoint": "localhost:9201",
    "log": {
      "level": "error"
    }
  },
  "client": {
    "log": {
      "level": "error"
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
var el, _ = elastic.NewElastic()
var searcher, _ = search.New()

func main() {
	// with database
	CleanDatabase()
	FillDatatabase()
	<-time.After(5 * time.Second)
	SearchFromDatabase()
	CleanDatabase()

	// with elastic
	CleanElastic()
	FillElastic()
	<-time.After(5 * time.Second)
	SearchFromElastic()
	CleanElastic()
}

func SearchFromDatabase() {

	result, err := searcher.NewDatabaseSearch(
		db.Select("*").
			From("public.person").
			OrderAsc("id_person")).
		Query(map[string]string{"first_name": "joao", "last_name": "ribeiro"}).
		Filters("first_name", "last_name").
		SearchFilters("first_name", "last_name").
		Search("joao").
		Bind(&[]Person{}).
		Path("http://teste.pt").
		Page(1).
		Size(3).
		MaxSize(10).
		Metadata("my-meta",
			db.Select("*").
				From("public.person").
				OrderAsc("id_person"),
			&[]Person{}).
		MetadataFunction("my-function", myDatabaseMetadataFunction, &[]Person{}).
		Fallback(searcher.NewElasticSearch(
			el.Search().
				Index("persons").
				Type("person")).
			Query(map[string]string{"first_name": "joao", "last_name": "ribeiro"}).
			Filters("first_name", "last_name").
			SearchFilters("first_name", "last_name").
			Search("joao").
			Bind(&[]Person{}).
			Path("http://teste.pt").
			Page(1).
			Size(3).
			MaxSize(10).
			Metadata("my-meta",
				el.Search().
					Index("persons").
					Type("person"),
				&[]Person{}).
			MetadataFunction("my-function", myElasticMetadataFunction, &[]Person{})).
		Exec()

	if err != nil {
		panic(err)
	}

	if result != nil {
		b, _ := json.MarshalIndent(result, "", "\t")
		fmt.Printf("\n\nSearch: %s", string(b))
	}
}

func SearchFromElastic() {

	result, err := searcher.NewElasticSearch(el.Search().
		Index("persons").
		Type("person")).
		Query(map[string]string{"first_name": "joao", "last_name": "ribeiro"}).
		Filters("first_name", "last_name").
		SearchFilters("first_name", "last_name").
		Search("joao").
		Bind(&[]Person{}).
		Path("http://teste.pt").
		Page(1).
		Size(3).
		MaxSize(10).
		Metadata("my-meta",
			el.Search().
				Index("persons").
				Type("person"),
			&[]Person{}).
		MetadataFunction("my-function", myElasticMetadataFunction, &[]Person{}).
		Fallback(searcher.NewDatabaseSearch(
			db.Select("*").
				From("public.person").
				OrderAsc("id_person")).
			Query(map[string]string{"first_name": "joao", "last_name": "ribeiro"}).
			Filters("first_name", "last_name").
			SearchFilters("first_name", "last_name").
			Search("joao").
			Bind(&[]Person{}).
			Path("http://teste.pt").
			Page(1).
			Size(3).
			MaxSize(10).
			Metadata("my-meta",
				db.Select("*").
					From("public.person").
					OrderAsc("id_person"),
				&[]Person{}).
			MetadataFunction("my-function", myDatabaseMetadataFunction, &[]Person{})).
		Exec()

	if err != nil {
		panic(err)
	}

	if result != nil {
		b, _ := json.MarshalIndent(result, "", "\t")
		fmt.Printf("\n\nSearch: %s", string(b))
	}
}

func myDatabaseMetadataFunction(result interface{}, object interface{}, metadata map[string]*search.Metadata) error {
	if result != nil {
		if persons, ok := result.([]Person); ok && len(persons) > 0 {
			_, err := db.Select("*").
				From("public.person").
				Where("id_person = ?", persons[0].IdPerson).
				OrderAsc("id_person").
				Load(object)
			return err
		}
	}
	return nil
}

func myElasticMetadataFunction(result interface{}, object interface{}, metadata map[string]*search.Metadata) error {
	if result != nil {
		if persons, ok := result.([]Person); ok && len(persons) > 0 {
			_, err := el.Search().
				Index("persons").
				Type("person").
				Object(object).Search()
			return err
		}
	}
	return nil
}

func FillDatatabase() {
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

func CleanDatabase() {
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

func FillElastic() {
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

		// document create with id
		response, err := el.Document().Index("persons").Type("person").Id(strconv.Itoa(i)).Body(person).Create()

		if err != nil {
			panic(err)
		} else {
			fmt.Printf("\ncreated a new person with id %s\n", response.ID)
		}
	}
	fmt.Printf("\nINSERTED")
}

func CleanElastic() {
	fmt.Println("\n\n:: DELETE")

	response, err := el.Index().Index("persons").Delete()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleted persons index ok: %t\n", response.Acknowledged)
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
		"my-function": [
			{
				"id_person": 1,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 1,
				"active": false,
				"fk_address": 1
			}
		],
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

:: DELETE

DELETED

:: DELETE
[IN] http client send Method[DELETE] Url[localhost:9201/persons] on Start[2019-01-29 01:47:47.686695 +0000 WET m=+5.150227440]
[OUT] http client send Method[DELETE] Url[localhost:9201/persons] on Start[2019-01-29 01:47:47.686695 +0000 WET m=+5.150227440] Elapsed[332.642934ms]

deleted persons index ok: true

DELETED

:: INSERT
[IN] http client send Method[POST] Url[localhost:9201/persons/person/1] on Start[2019-01-29 01:47:48.021699 +0000 WET m=+5.485230032]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/1] on Start[2019-01-29 01:47:48.021699 +0000 WET m=+5.485230032] Elapsed[777.68454ms]

created a new person with id 1
[IN] http client send Method[POST] Url[localhost:9201/persons/person/2] on Start[2019-01-29 01:47:48.799626 +0000 WET m=+6.263153169]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/2] on Start[2019-01-29 01:47:48.799626 +0000 WET m=+6.263153169] Elapsed[27.698774ms]

created a new person with id 2
[IN] http client send Method[POST] Url[localhost:9201/persons/person/3] on Start[2019-01-29 01:47:48.827593 +0000 WET m=+6.291119907]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/3] on Start[2019-01-29 01:47:48.827593 +0000 WET m=+6.291119907] Elapsed[31.019061ms]

created a new person with id 3
[IN] http client send Method[POST] Url[localhost:9201/persons/person/4] on Start[2019-01-29 01:47:48.858817 +0000 WET m=+6.322343090]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/4] on Start[2019-01-29 01:47:48.858817 +0000 WET m=+6.322343090] Elapsed[24.027461ms]

created a new person with id 4
[IN] http client send Method[POST] Url[localhost:9201/persons/person/5] on Start[2019-01-29 01:47:48.883 +0000 WET m=+6.346526084]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/5] on Start[2019-01-29 01:47:48.883 +0000 WET m=+6.346526084] Elapsed[28.970813ms]

created a new person with id 5
[IN] http client send Method[POST] Url[localhost:9201/persons/person/6] on Start[2019-01-29 01:47:48.912141 +0000 WET m=+6.375666906]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/6] on Start[2019-01-29 01:47:48.912141 +0000 WET m=+6.375666906] Elapsed[15.377101ms]

created a new person with id 6
[IN] http client send Method[POST] Url[localhost:9201/persons/person/7] on Start[2019-01-29 01:47:48.927688 +0000 WET m=+6.391213673]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/7] on Start[2019-01-29 01:47:48.927688 +0000 WET m=+6.391213673] Elapsed[14.819912ms]

created a new person with id 7
[IN] http client send Method[POST] Url[localhost:9201/persons/person/8] on Start[2019-01-29 01:47:48.942796 +0000 WET m=+6.406322005]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/8] on Start[2019-01-29 01:47:48.942796 +0000 WET m=+6.406322005] Elapsed[19.8029ms]

created a new person with id 8
[IN] http client send Method[POST] Url[localhost:9201/persons/person/9] on Start[2019-01-29 01:47:48.96281 +0000 WET m=+6.426335770]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/9] on Start[2019-01-29 01:47:48.96281 +0000 WET m=+6.426335770] Elapsed[13.96723ms]

created a new person with id 9
[IN] http client send Method[POST] Url[localhost:9201/persons/person/10] on Start[2019-01-29 01:47:48.977 +0000 WET m=+6.440526256]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/10] on Start[2019-01-29 01:47:48.977 +0000 WET m=+6.440526256] Elapsed[20.478616ms]

created a new person with id 10
[IN] http client send Method[POST] Url[localhost:9201/persons/person/11] on Start[2019-01-29 01:47:48.997635 +0000 WET m=+6.461160949]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/11] on Start[2019-01-29 01:47:48.997635 +0000 WET m=+6.461160949] Elapsed[23.854232ms]

created a new person with id 11
[IN] http client send Method[POST] Url[localhost:9201/persons/person/12] on Start[2019-01-29 01:47:49.021736 +0000 WET m=+6.485261214]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/12] on Start[2019-01-29 01:47:49.021736 +0000 WET m=+6.485261214] Elapsed[30.152495ms]

created a new person with id 12
[IN] http client send Method[POST] Url[localhost:9201/persons/person/13] on Start[2019-01-29 01:47:49.052129 +0000 WET m=+6.515654027]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/13] on Start[2019-01-29 01:47:49.052129 +0000 WET m=+6.515654027] Elapsed[37.621388ms]

created a new person with id 13
[IN] http client send Method[POST] Url[localhost:9201/persons/person/14] on Start[2019-01-29 01:47:49.090377 +0000 WET m=+6.553902418]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/14] on Start[2019-01-29 01:47:49.090377 +0000 WET m=+6.553902418] Elapsed[43.181175ms]

created a new person with id 14
[IN] http client send Method[POST] Url[localhost:9201/persons/person/15] on Start[2019-01-29 01:47:49.133724 +0000 WET m=+6.597248935]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/15] on Start[2019-01-29 01:47:49.133724 +0000 WET m=+6.597248935] Elapsed[28.277231ms]

created a new person with id 15
[IN] http client send Method[POST] Url[localhost:9201/persons/person/16] on Start[2019-01-29 01:47:49.162164 +0000 WET m=+6.625688878]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/16] on Start[2019-01-29 01:47:49.162164 +0000 WET m=+6.625688878] Elapsed[23.353889ms]

created a new person with id 16
[IN] http client send Method[POST] Url[localhost:9201/persons/person/17] on Start[2019-01-29 01:47:49.185707 +0000 WET m=+6.649231433]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/17] on Start[2019-01-29 01:47:49.185707 +0000 WET m=+6.649231433] Elapsed[22.273157ms]

created a new person with id 17
[IN] http client send Method[POST] Url[localhost:9201/persons/person/18] on Start[2019-01-29 01:47:49.208175 +0000 WET m=+6.671699144]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/18] on Start[2019-01-29 01:47:49.208175 +0000 WET m=+6.671699144] Elapsed[32.801504ms]

created a new person with id 18
[IN] http client send Method[POST] Url[localhost:9201/persons/person/19] on Start[2019-01-29 01:47:49.241422 +0000 WET m=+6.704946537]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/19] on Start[2019-01-29 01:47:49.241422 +0000 WET m=+6.704946537] Elapsed[45.687112ms]

created a new person with id 19
[IN] http client send Method[POST] Url[localhost:9201/persons/person/20] on Start[2019-01-29 01:47:49.287351 +0000 WET m=+6.750875240]
[OUT] http client send Method[POST] Url[localhost:9201/persons/person/20] on Start[2019-01-29 01:47:49.287351 +0000 WET m=+6.750875240] Elapsed[40.078909ms]

created a new person with id 20

[IN] http client send Method[POST] Url[localhost:9201/persons/_count] on Start[2019-01-29 01:47:54.328857 +0000 WET m=+11.792352950]
[OUT] http client send Method[POST] Url[localhost:9201/persons/_count] on Start[2019-01-29 01:47:54.328857 +0000 WET m=+11.792352950] Elapsed[19.266367ms]
[IN] http client send Method[POST] Url[localhost:9201/persons/_search?size=3&from=0] on Start[2019-01-29 01:47:54.34839 +0000 WET m=+11.811886152]
[OUT] http client send Method[POST] Url[localhost:9201/persons/_search?size=3&from=0] on Start[2019-01-29 01:47:54.34839 +0000 WET m=+11.811886152] Elapsed[16.904429ms]
[IN] http client send Method[GET] Url[localhost:9201/persons/_search] on Start[2019-01-29 01:47:54.365622 +0000 WET m=+11.829118651]
[OUT] http client send Method[GET] Url[localhost:9201/persons/_search] on Start[2019-01-29 01:47:54.365622 +0000 WET m=+11.829118651] Elapsed[15.296439ms]
[IN] http client send Method[GET] Url[localhost:9201/persons/_search] on Start[2019-01-29 01:47:54.381199 +0000 WET m=+11.844694617]
[OUT] http client send Method[GET] Url[localhost:9201/persons/_search] on Start[2019-01-29 01:47:54.381199 +0000 WET m=+11.844694617] Elapsed[12.052136ms]
[IN] http client send Method[POST] Url[localhost:9201/persons/_search?size=3&from=0] on Start[2019-01-29 01:47:54.393709 +0000 WET m=+11.857204913]
[OUT] http client send Method[POST] Url[localhost:9201/persons/_search?size=3&from=0] on Start[2019-01-29 01:47:54.393709 +0000 WET m=+11.857204913] Elapsed[33.738833ms]


Search: {
	"result": [
		{
			"id_person": 14,
			"first_name": "joao",
			"last_name": "ribeiro",
			"age": 14,
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
			"id_person": 3,
			"first_name": "joao",
			"last_name": "ribeiro",
			"age": 3,
			"active": false,
			"fk_address": 1
		}
	],
	"metadata": {
		"my-function": [
			{
				"id_person": 14,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 14,
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
				"id_person": 5,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 5,
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
				"id_person": 12,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 12,
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
				"id_person": 4,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 4,
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
			}
		],
		"my-meta": [
			{
				"id_person": 14,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 14,
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
				"id_person": 5,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 5,
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
				"id_person": 12,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 12,
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
				"id_person": 4,
				"first_name": "joao",
				"last_name": "ribeiro",
				"age": 4,
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

:: DELETE
[IN] http client send Method[DELETE] Url[localhost:9201/persons] on Start[2019-01-29 01:47:54.427835 +0000 WET m=+11.891330396]
[OUT] http client send Method[DELETE] Url[localhost:9201/persons] on Start[2019-01-29 01:47:54.427835 +0000 WET m=+11.891330396] Elapsed[383.538141ms]

deleted persons index ok: true

DELETED
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
