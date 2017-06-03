+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "New"
description = ""

[menu.main]
parent = "API"
identifier = "new"
weight = 5
+++

Instantiate an elasticsearch client with default values like so:

```go
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{})
}
```

By default the client will expect a running elasticsearch instance at localhost:9200. If your service is on a different
host or port you can pass it a url directly.

```go 
package main 
 
import (
    "github.com/b3ntly/elasticsearch"
)

func main(){
        client, err := elasticsearch.New(&elasticsearch.Options{ URL: "http://elasticsearch:9201"})
}
```

