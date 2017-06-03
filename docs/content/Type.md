+++
date = "2017-06-02T13:33:55-07:00"
description = ""
title = "Type"
draft = false

[menu.main]
identifier = "Type"
parent = "API"
weight = 30
+++

See [this link]("/jargon) for the purpose and definition of an elasticsearch Type. When a reference to a Type is
obtained, as in the below code, a type is not actually created in elasticsearch. The default behavior for elasticsearch
is to create the type as soon as a operation is performed that references it, i.e. an Insert or Update call, thus there
is no need to create it immediately.

```go
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{})
        
        // type is a reserved word in golang so we use the identifier collection
        collection := client.I("test").T("test")
}
```