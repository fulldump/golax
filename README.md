<img src="logo.png">

<p align="center">
<img src="https://api.travis-ci.org/fulldump/golax.svg?branch=master">
</p>

Golax is the official go implementation for the _Lax_ framework.

<!-- MarkdownTOC autolink=true bracket=round depth=4 -->

- [About Lax](#about-lax)
- [Getting started](#getting-started)
- [How interceptor works](#how-interceptor-works)
- [Handling parameters](#handling-parameters)
- [Sample use cases](#sample-use-cases)

<!-- /MarkdownTOC -->

Related docs:

* [Developer notes](doc/developer.md)
* [TODO list](doc/todo.md)

## About Lax

Lax wants to be the best _"user experience"_ for developers making REST APIs.

The design principles for _Lax_ are:

* The lowest language overhead
* Extremely fast to develop
* Very easy to read and trace.


## Getting started

```go
my_api := golax.NewApi()

my_api.Root.
    Interceptor(golax.InterceptorError).
    Interceptor(myLogingInterceptor)

my_api.Root.Node("hello").
    Method("GET", func(c *golax.Context) {
        // At this point, Root interceptors has been already executed
        fmt.Fprintln(c.Response, "Hello world!")
    })

my_api.Serve()
```

## How interceptor works

If I want to handle a `GET /users/1234/stats` request, all interceptors in nodes from `<root>` to `.../stats` are executed:

![Normal flow](figure_1_normal_flow.png)

To abort the execution, call to `c.Error(404, "Resource not found")`:

![Break flow](figure_2_break_flow.png)

## Handling parameters

```go
my_api := golax.NewApi()

my_api.Root.
    Node("users").
    Node("{user_id}").
    Method("GET", func (c *golax.Context) {
        fmt.Fprintln(c.Response, "You are looking for user " + c.Parameter)
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
