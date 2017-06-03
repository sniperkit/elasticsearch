+++
date = "2017-06-02T13:25:49-07:00"
draft = false
title = "Dependencies"
description = ""

[menu.main]
parent = "Architecture"
identifier = "Dependencies"
weight = 10
+++

Elasticsearch has two minimal dependencies github.com/stretchr/testify and github.com/hashicorp/go-cleanhttp which 
are vendored
using the Glide vendoring library. Aside from that our test harness requires an actively running elasticsearch service
found at localhost:9200. 

In future versions this library will contain a fully mocked elasticsearch service that removes elasticsearch as a
dependency for the test suite.