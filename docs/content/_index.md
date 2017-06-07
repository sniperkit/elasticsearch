+++
draft = false
title = "Home"
description = ""
date = "2017-04-24T18:36:24+02:00"
+++

<span id="sidebar-toggle-span">
<a href="#" id="sidebar-toggle" data-sidebar-toggle="">Home<i class="fa fa-bars"></i></a>
<i class='fa fa-github'></i> 
</span>


# Abstract

Elasticsearch is a distributed full-text search engine built on Lucene. This library creates a Mongo DB like 
interface so Elasticsearch may be used as a standalone database for small web applications.

Note the DBMS and this library are both called 'elasticsearch'. As often as possible I will use the term 'ES' to 
signify elasticsearch the database.

### Provides

* A transparent protocol for interacting with the Elasticsearch REST API
* Interfaces based on the the mgo library for Mongo DB.
* Solid test harness for future extensibility
* Bulk API support
* elasticsearch/mock - A functioning HTTP server which replicated the base functionality of Elasticsearch
* Roadmap for deep querying support