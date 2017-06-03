+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "DeleteById"
description = ""

[menu.main]
parent = "Type"
identifier = "Delete"
weight = 50
+++

Delete a document by ID.

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
        err := collection.DeleteById(doc.ID)
}
```