+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "Drop"
description = ""

[menu.main]
parent = "Index"
identifier = "Drop"
weight = 60
+++

Drop an index and all its underlying documents.

```go
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{})
        err := client.I("test").Drop()
}
```