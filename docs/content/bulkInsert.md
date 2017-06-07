+++
date = "2017-06-02T13:30:25-07:00"
description = ""
title = "BulkInsert"

[menu.main]
identifier = "BulkInsertInsert"
parent = "Type"
weight = 11
+++

Insert multiple JSON objects in the form of byte slice(s). If an _id property is specified it will be set as the
_id of the document. If one is not provided the _id will be automatically
generated.

If you do specify an _id and the _id already exists the API call will return
an error.

Notably the ES Bulk API is not transactional and therefore may only partially complete
bulk request.

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
        
        // inserted represents a slice of IDs
        // this call will return an error of not all documents were inserted.
        inserted, err := collection.BulkInsert([][]byte{[]byte("{\"message\": \"hello, world\"}")})
}
```