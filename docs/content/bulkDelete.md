+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "BulkDelete"
description = ""

[menu.main]
parent = "Type"
identifier = "BulkDelete"
weight = 51
+++

Bulk delete a document by their ID(s).

The ES Bulk API is not transactional and may only partially complete a series of deletions in a bulk request.

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
        
        // deleted is a slice of strings representing the deleted documents
        deleted, err := collection.BulkDelete(ID)
}
```