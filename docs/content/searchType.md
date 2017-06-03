+++
date = "2017-06-02T13:30:25-07:00"
description = ""
title = "Search"

[menu.main]
identifier = "SearchType"
parent = "Type"
weight = 20
+++

Search returns a slice of Documents representing the search results of your query. Indexes and Types both contain a 
search method that will search their respective domain, i.e. searching an index will search all documents in an index
while searching a type will return all documents in the index that also fulfill the given type.

Currently search only accepts a simple querystring parameter that can be used to test for exact string matches and
property matches.

```go
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
    "os"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{})
        
        if err != nil {
                os.Exit(-1)
        }
        
        collection := client.I("test").T("test")
        doc, err := collection.Insert([]byte("{\"message\": \"hello, world\"}"))
        docs, err := collection.Search("hello")
}
```