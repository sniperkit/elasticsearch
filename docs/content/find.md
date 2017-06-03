+++
date = "2017-06-02T13:30:45-07:00"
description = ""
title = "Find"

[menu]

  [menu.main]
    identifier = "Find"
    parent = "Type"
    weight = 25

+++

Find() implements search under the hood and exists purely to satisfy a mgo like interface.

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
        docs, err := collection.Find("hello")
        docs, err := collection.Find("message:hello, world")
}
```
