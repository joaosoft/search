# elastic
[![Build Status](https://travis-ci.org/joaosoft/elastic.svg?branch=master)](https://travis-ci.org/joaosoft/elastic) | [![codecov](https://codecov.io/gh/joaosoft/elastic/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/elastic) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/elastic)](https://goreportcard.com/report/github.com/joaosoft/elastic) | [![GoDoc](https://godoc.org/github.com/joaosoft/elastic?status.svg)](https://godoc.org/github.com/joaosoft/elastic)

A simple and fast elastic-search client.

## Support for 
* Index (Create / Exists / Delete)  [with or without mapping]
* Document (Create / Update / Delete)
* Search
* Bulk (Index / Create / Update / Delete)

* The search can be done with a template to be faster than other complicated frameworks.

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/elastic
```

## Usage 
This examples are available in the project at [elastic/examples](https://github.com/joaosoft/elastic/tree/master/examples)

### Templates
#### get.example.1.template
```
{
  "query": {
    "bool": {
      "must": {
        "term": {
          {{ range $key, $value := .Data }}
             "{{ $key }}": "{{ $value }}"
             {{ if (gt (len $.Data) 1) }}
                 ,
             {{ end }}
          {{ end }}
        }
      }
    }
  },
  "sort": [
    {
      "age": {
        "order": "desc"
      }
    }
  ]

  {{ if (gt $.From 0) }}
    ,
    "from": {{.From}}
  {{ end }}

  {{ if (gt $.Size 0) }}
    ,
  " size": {{.Size}}
  {{ end }}
}
```

>### Implementation
```go
// create a client
type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

//var client = elastic.NewElastic()
// you can define the configuration without having a configuration file
var config, _, _ = elastic.NewConfig("localhost:9201")
var client, _ = elastic.NewElastic(elastic.WithConfiguration(config.Elastic))

func main() {

	// delete indexes
	fmt.Println(":: DELETE INDEX")
	deleteIndex()

	// index exists
	fmt.Println(":: EXISTS INDEXES ?")
	existsIndex("persons")
	existsIndex("bananas") // don't exist

	// index create with mapping
	fmt.Println(":: CREATE INDEX WITH MAPPING")
	createIndexWithMapping()

	// document create
	fmt.Println(":: CREATE DOCUMENT 1")
	createDocumentWithId("1")
	fmt.Println(":: CREATE DOCUMENT 2")
	createDocumentWithId("2")
	fmt.Println(":: CREATE DOCUMENT WITHOUT ID")
	generatedId := createDocumentWithoutId()

	// document update
	fmt.Println(":: UPDATE DOCUMENT 1")
	updateDocumentWithId("1")
	fmt.Println(":: UPDATE DOCUMENT 2")
	updateDocumentWithId("2")

	// document search
	// wait elastic to index the last update...
	fmt.Println(":: SEARCH DOCUMENT  WITH 'luis'")
	<-time.After(time.Second * 2)
	searchDocument("luis")

	// count index documents
	fmt.Println(":: COUNT DOCUMENTS ON INDEX WITH 'luis'")
	countOnIndex("luis")
	fmt.Println(":: COUNT DOCUMENTS ON DOCUMENT WITH 'luis'")
	countOnDocument("luis")

	// document delete
	fmt.Println(":: COUNT DOCUMENT WITH GENERATED ID")
	deleteDocumentWithId("1")
	deleteDocumentWithId("2")
	deleteDocumentWithId(generatedId)

	// bulk
	fmt.Println(":: BULK OPERATIONS")
	bulkCreate()
	bulkIndex()
	bulkDelete()

	// queue bulk create
	fmt.Println(":: QUEUE BULK CREATE")
	queueBulkCreate()

	// index delete
	fmt.Println(":: DELETE INDEX")
	deleteIndex()

}

func createIndexWithMapping() {
	// you can define the configuration without having a configuration file
	//client1 := elastic.NewElastic(elastic.WithConfiguration(elastic.NewConfig("http://localhost:9200")))

	response, err := client.Index().Index("persons").Body([]byte(`
{
  "mappings": {
    "person": {
      "properties": {
        "age": {
          "type": "long"
        },
        "name": {
          "type": "text",
          "fields": {
            "keyword": {
              "type": "keyword",
              "ignore_above": 256
            }
          }
        }
      }
    }
  }
}
`)).Create()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ncreated mapping for persons index ok: %t\n", response.Acknowledged)
	}
}

func deleteIndex() {

	response, err := client.Index().Index("persons").Delete()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleted persons index ok: %t\n", response.Acknowledged)
	}
}

func existsIndex(index string) {

	exists, err := client.Index().Index(index).Exists()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nexists index? %t\n", exists)
	}
}

func countOnIndex(name string) int64 {

	d1 := elastic.CountTemplate{Data: map[string]interface{}{"name": name}}

	// index count
	dir, _ := os.Getwd()
	response, err := client.Search().
		Index("persons").
		Template(dir+"/examples/templates", "get.example.count.template", &d1, false).
		Count()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ncount persons with name %s: %d\n", name, response.Count)
	}

	return response.Count
}

func countOnDocument(name string) int64 {

	d1 := elastic.CountTemplate{Data: map[string]interface{}{"name": name}}

	// index count
	dir, _ := os.Getwd()
	response, err := client.Search().
		Index("persons").
		Type("person").
		Template(dir+"/examples/templates", "get.example.count.template", &d1, false).
		Count()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ncount persons with name %s: %d\n", name, response.Count)
	}

	return response.Count
}

func createDocumentWithId(id string) {

	// document create with id
	age, _ := strconv.Atoi(id)
	response, err := client.Document().Index("persons").Type("person").Id(id).Body(person{
		Name: "joao",
		Age:  age + 20,
	}).Create()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ncreated a new person with id %s\n", response.ID)
	}
}

func createDocumentWithoutId() string {

	// document create without id
	response, err := client.Document().Index("persons").Type("person").Body(person{
		Name: "joao",
		Age:  30,
	}).Create()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ncreated a new person with id %s\n", response.ID)
	}

	return response.ID
}

func updateDocumentWithId(id string) {

	// document update with id
	age, _ := strconv.Atoi(id)
	response, err := client.Document().Index("persons").Type("person").Id(id).Body(person{
		Name: "luis",
		Age:  age + 20,
	}).Update()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nupdated person with id %s\n", response.ID)
	}
}

func deleteDocumentWithId(id string) {

	response, err := client.Document().Index("persons").Type("person").Id(id).Delete()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleted person with id %s\n", response.ID)
	}
}

func searchDocument(name string) {
	var data []person

	d1 := elastic.SearchTemplate{Data: map[string]interface{}{"name": name}}

	// document search
	dir, _ := os.Getwd()
	_, err := client.Search().
		Index("persons").
		Type("person").
		Object(&data).
		Template(dir+"/examples/templates", "get.example.search.template", &d1, false).
		Search()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nsearch person by name:%s %+v\n", name, data)
	}
}

func bulkCreate() {
	bulk := client.Bulk()

	// document create with id
	id := "1"
	err := bulk.Index("persons").Type("person").Id(id).Body(person{
		Name: "joao",
		Age:  1,
	}).DoCreate()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nadding a new person with id %s\n", id)
	}

	id = "2"
	err = bulk.Index("persons").Type("person").Id(id).Body(person{
		Name: "tiago",
		Age:  2,
	}).DoCreate()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nadding a new person with id %s\n", id)
	}

	fmt.Println("executing bulk")
	_, err = bulk.Execute()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("success!")
	}
}

func bulkIndex() {
	bulk := client.Bulk()

	// document create with id
	id := "3"
	err := bulk.Index("persons").Type("person").Id(id).Body(person{
		Name: "joao",
		Age:  1,
	}).DoIndex()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nadding a new person with id %s\n", id)
	}

	id = "4"
	err = bulk.Index("persons").Type("person").Id(id).Body(person{
		Name: "tiago",
		Age:  2,
	}).DoCreate()

	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\nadding a new person with id %s\n", id)
	}

	fmt.Println("executing bulk")
	_, err = bulk.Execute()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("success!")
	}
}

func bulkDelete() {
	bulk := client.Bulk()

	// document delete with id
	id := "1"
	err := bulk.Index("persons").Type("person").Id(id).DoDelete()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleting the person with id %s\n", id)
	}

	id = "2"
	err = bulk.Index("persons").Type("person").Id(id).DoDelete()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleting the person with id %s\n", id)
	}

	// document delete with id
	id = "3"
	err = bulk.Index("persons").Type("person").Id(id).DoDelete()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleting the person with id %s\n", id)
	}

	id = "4"
	err = bulk.Index("persons").Type("person").Id(id).DoDelete()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("\ndeleting the person with id %s\n", id)
	}

	fmt.Println("executing bulk")
	_, err = bulk.Execute()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("success!")
	}
}

func queueBulkCreate() {
	// create process manager
	pm := manager.NewManager(manager.WithRunInBackground(false))

	// create queue
	bulkWorkqueueConfig := manager.NewBulkWorkListConfig("queue_001", 100, 10, 2, time.Second*2, manager.FIFO)
	bulkWorkqueue := pm.NewSimpleBulkWorkList(bulkWorkqueueConfig, bulkWorkHandler, bulkWorkRecoverHandler, bulkWorkRecoverWastedRetriesHandler)
	pm.AddWorkList("bulk_queue", bulkWorkqueue)

	if err := bulkWorkqueue.Start(); err != nil {
		log.Errorf("MAIN: error starting bulk workqueue %s", err)
	}

	// add job to queue
	go func() {
		nJobs := 20000
		for i := 1; i <= nJobs; i++ {
			bulkWorkqueue.AddWork(strconv.Itoa(i),
				&person{
					Name: fmt.Sprintf("name %d", i),
					Age:  i,
				})
		}
	}()

	<-time.After(30 * time.Second)
}

func bulkWorkHandler(works []*manager.Work) error {
	log.Infof("handling works with length %d!", len(works))

	bulk := client.Bulk()

	// handle works on elastic bulk
	var err error
	for _, work := range works {
		if err = bulk.Index("persons").Type("person").Id(work.Id).Body(work.Data).DoCreate(); err != nil {
			panic(err)
			return err
		}
	}

	_, err = bulk.Execute()
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("success!")
	}

	return nil
}

func bulkWorkRecoverHandler(list manager.IList) error {
	fmt.Printf("\nrecovering list with length %d", list.Size())
	return nil
}

func bulkWorkRecoverWastedRetriesHandler(id string, data interface{}) error {
	fmt.Printf("\nrecovering work with id: %s, data: %+v", id, data)
	return nil
}
```

>### Result
```
:: DELETE INDEX
Status[404] Method[DELETE] Url[/persons] on Start[2019-09-13T13:58:32+01:00] Elapsed[21.235791ms]

deleted persons index ok: false
:: EXISTS INDEXES ?
Status[404] Method[HEAD] Url[/persons] on Start[2019-09-13T13:58:33+01:00] Elapsed[11.317321ms]

exists index? true
Status[404] Method[HEAD] Url[/bananas] on Start[2019-09-13T13:58:33+01:00] Elapsed[15.069151ms]

exists index? true
:: CREATE INDEX WITH MAPPING
Status[200] Method[PUT] Url[/persons] on Start[2019-09-13T13:58:33+01:00] Elapsed[4.415234955s]

created mapping for persons index ok: true
:: CREATE DOCUMENT 1
Status[201] Method[POST] Url[/persons/person/1] on Start[2019-09-13T13:58:37+01:00] Elapsed[116.445175ms]

created a new person with id 1
:: CREATE DOCUMENT 2
Status[201] Method[POST] Url[/persons/person/2] on Start[2019-09-13T13:58:37+01:00] Elapsed[48.192381ms]

created a new person with id 2
:: CREATE DOCUMENT WITHOUT ID
Status[201] Method[POST] Url[/persons/person] on Start[2019-09-13T13:58:37+01:00] Elapsed[112.728693ms]

created a new person with id sQS0Km0BmsVYzvqSO8wY
:: UPDATE DOCUMENT 1
Status[200] Method[PUT] Url[/persons/person/1] on Start[2019-09-13T13:58:37+01:00] Elapsed[26.134477ms]

updated person with id 1
:: UPDATE DOCUMENT 2
Status[200] Method[PUT] Url[/persons/person/2] on Start[2019-09-13T13:58:37+01:00] Elapsed[49.6229ms]

updated person with id 2
:: SEARCH DOCUMENT  WITH 'luis'
Status[200] Method[POST] Url[/persons/_search] on Start[2019-09-13T13:58:39+01:00] Elapsed[39.807747ms]

search person by name:luis []
:: COUNT DOCUMENTS ON INDEX WITH 'luis'
Status[200] Method[POST] Url[/persons/_count] on Start[2019-09-13T13:58:39+01:00] Elapsed[275.765858ms]

count persons with name luis: 0
:: COUNT DOCUMENTS ON DOCUMENT WITH 'luis'
Status[200] Method[POST] Url[/persons/_count] on Start[2019-09-13T13:58:40+01:00] Elapsed[170.410779ms]

count persons with name luis: 0
:: COUNT DOCUMENT WITH GENERATED ID
Status[200] Method[DELETE] Url[/persons/person/1] on Start[2019-09-13T13:58:40+01:00] Elapsed[173.527989ms]

deleted person with id 1
Status[200] Method[DELETE] Url[/persons/person/2] on Start[2019-09-13T13:58:40+01:00] Elapsed[59.927572ms]

deleted person with id 2
Status[200] Method[DELETE] Url[/persons/person/sQS0Km0BmsVYzvqSO8wY] on Start[2019-09-13T13:58:40+01:00] Elapsed[367.634448ms]

deleted person with id sQS0Km0BmsVYzvqSO8wY
:: BULK OPERATIONS

adding a new person with id 1

adding a new person with id 2
executing bulk
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:40+01:00] Elapsed[98.144062ms]
success!
adding a new person with id 3

adding a new person with id 4
executing bulk
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:41+01:00] Elapsed[189.19494ms]
success!
deleting the person with id 1

deleting the person with id 2

deleting the person with id 3

deleting the person with id 4
executing bulk
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:41+01:00] Elapsed[50.71554ms]
success!:: QUEUE BULK CREATE
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:41:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:41+01:00] Elapsed[439.981927ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:41:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:41+01:00] Elapsed[198.282728ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:41:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:41+01:00] Elapsed[213.071722ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:42:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:42+01:00] Elapsed[247.426975ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:42:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:42+01:00] Elapsed[319.394449ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:42:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:42+01:00] Elapsed[377.465058ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:43:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[905.412147ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:44:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[1.147159558s]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[1.147512043s]
success!success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:44:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:44:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[1.613820621s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:44:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[1.870385532s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:45:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[2.661330737s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:45:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[3.720321356s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:46:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[4.18292822s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:47:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[4.210801216s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:47:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:43+01:00] Elapsed[4.255845561s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:47:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:44+01:00] Elapsed[3.932291525s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:48:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:44+01:00] Elapsed[4.159273553s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:48:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:44+01:00] Elapsed[4.609895408s]

recovering list with length 17165{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:49:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[0] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:44+01:00] Elapsed[5.001851929s]

recovering list with length 17065{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:49:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[0] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:45+01:00] Elapsed[5.00494695s]

recovering list with length 16965{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:50:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:45+01:00] Elapsed[4.763509759s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:50:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:46+01:00] Elapsed[3.97617553s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:50:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:47+01:00] Elapsed[3.517492575s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:50:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:47+01:00] Elapsed[3.580138183s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:47+01:00] Elapsed[3.547918936s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:48+01:00] Elapsed[3.413002589s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:48+01:00] Elapsed[3.075790415s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:49+01:00] Elapsed[2.715319313s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:50+01:00] Elapsed[1.748025446s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:49+01:00] Elapsed[2.115694509s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:51:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:50+01:00] Elapsed[1.958935965s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:52:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:50+01:00] Elapsed[2.050065457s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:52:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:50+01:00] Elapsed[1.907570135s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:52:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.077481421s]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.090619887s]
success!success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:53:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:53:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[1.952769005s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:53:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.175756422s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:53:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.15950508s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:53:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.323748995s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:54:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:51+01:00] Elapsed[2.625130992s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:54:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:52+01:00] Elapsed[2.070779945s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:54:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:52+01:00] Elapsed[2.059658107s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:54:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:52+01:00] Elapsed[2.03033224s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:54:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:53+01:00] Elapsed[2.075507055s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:55:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:53+01:00] Elapsed[2.439292764s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:55:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:53+01:00] Elapsed[1.811780147s]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:53+01:00] Elapsed[2.166005151s]
success!success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:55:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:55:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:53+01:00] Elapsed[1.959211779s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:55:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:54+01:00] Elapsed[1.807342011s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:54+01:00] Elapsed[1.46737623s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:54+01:00] Elapsed[1.564310936s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:54+01:00] Elapsed[1.510591798s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:54+01:00] Elapsed[1.719323571s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:55+01:00] Elapsed[1.588701082s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:55+01:00] Elapsed[1.278724736s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:56:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:55+01:00] Elapsed[1.718028078s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:55+01:00] Elapsed[1.819044919s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:55+01:00] Elapsed[1.75113316s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.689605882s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.565944471s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.789003032s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.610329072s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:57:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.436485117s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:58:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.480721966s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:58:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.389606351s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:58:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:56+01:00] Elapsed[1.88361716s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:58:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.418232025s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:58:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.407353734s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:59:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.466319058s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:59:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.5826686s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:59:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.928709837s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:59:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:57+01:00] Elapsed[1.929256988s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:58:59:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:58+01:00] Elapsed[1.973185976s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:58+01:00] Elapsed[1.873539511s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:58+01:00] Elapsed[1.54940218s]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:58+01:00] Elapsed[1.552183951s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:58+01:00] Elapsed[1.59217799s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:59+01:00] Elapsed[1.472128878s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:59+01:00] Elapsed[1.389092457s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:59+01:00] Elapsed[1.157879224s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:59+01:00] Elapsed[870.375509ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:58:59+01:00] Elapsed[774.044196ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[669.533882ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[590.004246ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[467.769253ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[532.764842ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[394.688819ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[411.011855ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[418.469861ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:00:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[434.895767ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[439.760596ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[442.70354ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[478.885684ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[449.657235ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[477.702121ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[501.51052ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[532.086787ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[565.763775ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:00+01:00] Elapsed[520.071ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[466.607987ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[496.990002ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[553.476716ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[527.756385ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:01:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[806.47811ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[821.273632ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[763.212143ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[745.978277ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[800.176633ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[802.756072ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[817.559681ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[752.515078ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[710.490962ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:01+01:00] Elapsed[783.261523ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[589.511546ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[660.007353ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[735.92878ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[738.52141ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[714.896762ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:02:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[849.547496ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[841.285181ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[988.560743ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[1.107840265s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[1.086501274s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[1.032504289s]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[949.798564ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[990.201515ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[973.675585ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:02+01:00] Elapsed[936.000033ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[829.649349ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:03:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[866.094659ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[715.861696ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[590.671932ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[547.382082ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[550.566764ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[537.839107ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[453.900829ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[450.121584ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[413.328177ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:03+01:00] Elapsed[396.69272ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[386.580576ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[372.129023ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[393.529319ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[425.049786ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[405.136084ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[396.330373ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[389.58638ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[375.077702ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[413.835554ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[416.4085ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[405.855636ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[433.857421ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[406.827459ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[418.832457ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:04:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[396.839596ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[417.477493ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[404.182341ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[420.858559ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[491.799554ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[401.799091ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[407.271956ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[395.158636ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[388.975689ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:04+01:00] Elapsed[408.431226ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[408.799113ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[416.244833ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[409.237111ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[362.591567ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[408.726396ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[402.752503ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[400.425795ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[412.932506ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[401.770802ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[428.138748ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[448.672388ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[497.56238ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:05:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[510.9651ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[523.196423ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[539.908112ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[547.990185ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[549.060294ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[574.914387ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[544.214314ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[475.681096ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[445.852649ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[407.923029ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[383.849324ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[370.408211ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[357.492516ms]
success!{"prefixes":{"level":"info","timestamp":"2019-09-13 13:59:06:19"},"message":"handling works with length [100]!","sufixes":{"ip":"192.168.1.30"}}
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[347.662007ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[317.083791ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[320.113222ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[307.105526ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[297.68023ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[313.438435ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:05+01:00] Elapsed[1.39032216s]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[316.021488ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[339.105443ms]
Status[200] Method[POST] Url[/_bulk] on Start[2019-09-13T13:59:06+01:00] Elapsed[319.386954ms]
success!:: DELETE INDEX
Status[200] Method[DELETE] Url[/persons] on Start[2019-09-13T13:59:11+01:00] Elapsed[438.821611ms]

deleted persons index ok: true
```

## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
