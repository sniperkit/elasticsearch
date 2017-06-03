+++
date = "2017-06-02T16:49:39-07:00"
description = ""
title = "FindById"

[menu]

  [menu.main]
    identifier = "FindById"
    parent = "Type"
    weight = 27

+++

Return a document by ID. Will return an error if the document doesn't exist.

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
        docs, err := collection.FindById(doc.ID)  
}
```
