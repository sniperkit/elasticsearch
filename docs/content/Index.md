+++
date = "2017-06-02T13:33:55-07:00"
description = ""
title = "Index"
draft = false
url = "/search"

[menu.main]
identifier = "Index"
parent = "API"
weight = 20
+++

See [this link]("/jargon) for the purpose and definition of an elasticsearch Index. When a reference to an Index is
obtained, as in the below code, an Index is not actually created in elasticsearch. The default behavior for elasticsearch
is to create the Index as soon as an operation is performed that references it, i.e. an Insert or Update call, thus there
is no need to create it immediately.

```go
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{})
        index := client.I("test")
}
```