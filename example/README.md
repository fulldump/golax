# Example project

The example project implements a sample REST API with `golax`.

<!-- MarkdownTOC autolink=true bracket=round depth=4 -->

- [Up and running](#up-and-running)
- [The API](#the-api)
    - [List user ids `GET /service/v1/users/`](#list-user-ids-get-servicev1users)
    - [Create new user `POST /service/v1/users/`](#create-new-user-post-servicev1users)
    - [Get a user `GET /service/v1/users/{user_id}`](#get-a-user-get-servicev1usersuser_id)
    - [Modify a user `POST /service/v1/users/{user_id}`](#modify-a-user-post-servicev1usersuser_id)
    - [Delete a user `DELETE /service/v1/users/{user_id}`](#delete-a-user-delete-servicev1usersuser_id)

<!-- /MarkdownTOC -->

## Up and running

How to build and run:

```sh
make example && ./_vendor/bin/example
```

## The API

It is a CRUD over a `users` collection.

### List user ids `GET /service/v1/users/`

Command:
```sh
curl -i http://localhost:8000/service/v1/users/
```

Result:
```http
HTTP/1.1 200 OK
Date: Sun, 31 Jan 2016 21:15:50 GMT
Content-Length: 8
Content-Type: text/plain; charset=utf-8

[1,2,3]
```

### Create new user `POST /service/v1/users/`

Command:
```sh
curl -i http://localhost:8000/service/v1/users/ --data '{"name":"Oscar"}'
```

Result:
```http
HTTP/1.1 201 Created
Date: Sun, 31 Jan 2016 21:32:38 GMT
Content-Length: 9
Content-Type: text/plain; charset=utf-8

{"id":4}
```

### Get a user `GET /service/v1/users/{user_id}`

Command:
```sh
curl -i http://localhost:8000/service/v1/users/4
```

Result:
```http
HTTP/1.1 200 OK
Date: Sun, 31 Jan 2016 21:33:36 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

{"name":"Oscar","age":0,"introduction":""}
```

### Modify a user `POST /service/v1/users/{user_id}`

Command:
```sh
curl -i http://localhost:8000/service/v1/users/4 \
--data '{"age":70, "introduction": "Hello, I like golax"}'
```

Result:
```http
HTTP/1.1 200 OK
Date: Sun, 31 Jan 2016 21:35:12 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8

```

### Delete a user `DELETE /service/v1/users/{user_id}`

Command:
```sh
curl -i -X DELETE http://localhost:8000/service/v1/users/1
```

Result:
```http
HTTP/1.1 200 OK
Date: Sun, 31 Jan 2016 21:36:40 GMT
Content-Length: 0
Content-Type: text/plain; charset=utf-8

```

