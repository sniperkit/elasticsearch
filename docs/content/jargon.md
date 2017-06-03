+++
date = "2017-06-02T13:23:58-07:00"
description = ""
title = "Jargon"
draft = false

[menu.main]
parent = ""
identifier = "Getting Started"
weight = 30

+++

#### Comparisons to other DBMS

* MySQL => Databases => Tables => Columns/Rows

* MongoDB => Databases => Collections > Documents
  
* Elasticsearch => Indices => Types => Documents with Properties

#### Comparisons to MGO library

mgo:
```golang
session, err := mgo.Dial(url)
c := session.DB(database).C(collection)
err := c.Find(query).One(&result)

```

elasticsearch:
```go 
client := elasticsearch.Client(&elasticSearch.Options{})
collection := client.I("test").C("test")
documents, err := client.Search("golang")
```