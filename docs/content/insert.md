+++
date = "2017-06-02T13:30:25-07:00"
description = ""
title = "Insert"

[menu.main]
identifier = "Insert"
parent = "Type"
weight = 10
+++

Insert a JSON object in the form of a byte slice as a new
document in elasticsearch. If an _id property is specified it will be set as the
_id of the document. If one is not provided the _id will be automatically
generated.

If you do specify an _id and the _id already exists the API call will return
an error.

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
        collection.Insert([]byte("{\"message\": \"hello, world\"}"))
}
```