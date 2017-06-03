+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "UpdateById"
description = ""

[menu.main]
parent = "Type"
identifier = "UpdateById"
weight = 40
+++

Update a document by ID. If the document does not exist it will be inserted, which
is the default behavior of elasticsearch.

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
        err = collection.UpdateById(doc.ID, []byte("{\"message\": \"bye, world\"}"))
}
```