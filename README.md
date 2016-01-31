<img src="logo.png">

<p align="center">
<img src="https://api.travis-ci.org/fulldump/golax.svg?branch=master">
</p>

Golax is the official go implementation for the _Lax_ pattern.

<!-- MarkdownTOC autolink=true bracket=round depth=4 -->

- [About Lax](#about-lax)
- [Getting started](#getting-started)
- [Sample use cases](#sample-use-cases)

<!-- /MarkdownTOC -->

## About Lax

Lax wants to be the best _"user experience"_ for developers.

The design principles for _Lax_ are:

* The lowest language overhead
* Extremely fast to develop
* Very easy to read and trace.


## Getting started

```go
my_api := golax.NewApi()

my_api.Root.Node("hello").
Method("GET", func(c *golax.Context) {
    fmt.Fprintln(c.Response, "Hello world!")
})

my_api.Serve()
```

## Sample use cases

TODO: put here some examples to cover cool things:

* parameters
* fluent implementation
* node cycling
* readability
* node preference
* sample logging middleware
* sample auth middleware
* sample api errors
