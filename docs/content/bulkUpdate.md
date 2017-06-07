+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "BulkUpdate"
description = ""

[menu.main]
parent = "Type"
identifier = "BulkUpdate"
weight = 41
+++

Update document(s) by ID. If the document does not exist it will be inserted, which
is the default behavior of elasticsearch.

Notably the ES Bulk API is not transactionally and may only partially complete a bulk request,
such as updating some but not all of the requested documents.

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
        ID, err := collection.Insert([]byte("{\"message\": \"hello, world\"}"))
       
        // updates is a slice of strings indicating the recently updated documents
        updates, err := collection.BulkUpdate([]*elasticsearch.Document{ &elasticsearch.Document{ ID: ID, Body: []byte("{\"message\": \"bye, world\"}"}})
}
```