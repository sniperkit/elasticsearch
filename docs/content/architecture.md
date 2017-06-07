+++
date = "2017-06-02T13:33:55-07:00"
description = ""
title = "Architecture"
draft = false

[menu.main]
identifier = "Architecture"
parent = ""
weight = 20

+++

{{<mermaid>}}
sequenceDiagram
    participant User
    participant Library
    User->>Library: String or Byte Slice
    participant Elasticsearch
    Library->>Elasticsearch: JSON POST
    Elasticsearch->>Library: JSON Response
    Library->>User: Document(s) or error
{{< /mermaid >}}

<br />
<br />

{{% alert theme="warning" %}}[options.go](https://github.com/b3ntly/elasticsearch/blob/master/options.go): Library configuration including default value initialization{{% /alert %}}

{{% alert theme="warning" %}}[client.go](https://github.com/b3ntly/elasticsearch/blob/master/client.go): Public API{{% /alert %}}

{{% alert theme="warning" %}}[rest.go](https://github.com/b3ntly/elasticsearch/blob/master/rest.go): REST interface with Elasticsearch{{% /alert %}}

{{% alert theme="warning" %}}[decode.go](https://github.com/b3ntly/elasticsearch/blob/master/decode.go): Map Elasticsearch responses to generic Document(s){{% /alert %}}

## Package Mock

Elasticsearch also provides a package called mock that contains a server implementation in Golang which replicates
the basic behavior of ES. Elasticsearch uses mock internally for added testing, though it should be noted
that it only very basically replicates ES (i.e. it is not sharded and represents data only in a nested map structure
protected by a RWMutex).