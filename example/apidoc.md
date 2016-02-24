# API Documentation
_Example_ is a demonstration REST API that implements a CRUD over a collection
of users, stored in memory.

All API calls:
* are returning errors with the same JSON format and
* are logging all request to standard output.

**Interceptors applied to all API:**  [`Log`](#interceptor-log)  [`Error`](#interceptor-error) 

## /service/v1/users

Resource users to list and create elements. It does not support pagination,
sorting or filtering.

**Interceptors chain:**  [`Log`](#interceptor-log)  [`Error`](#interceptor-error) 

**Methods:**  [`GET`](#get-servicev1users)  [`POST`](#post-servicev1users) 

### GET /service/v1/users

Return a list with a list of user ids:

```json
[1,2,3]
```

### POST /service/v1/users

Create a user:
```sh
curl http://localhost:8000/service/v1/users --data '{"name": "John"}'
```
And return the user id:
```json
{"id":4}
```

## /service/v1/users/{user_id}

Resource user to retrieve, modify and delete. A user has this structure:

```json
{
	"name": "Menganito Menganez",
	"age": 30,
	"introduction": "Hi, I like wheels and cars"
}
```

**Interceptors chain:**  [`Log`](#interceptor-log)  [`Error`](#interceptor-error)  [`User`](#interceptor-user) 

**Methods:**  [`GET`](#get-servicev1usersuser_id)  [`POST`](#post-servicev1usersuser_id)  [`DELETE`](#delete-servicev1usersuser_id) 

### GET /service/v1/users/{user_id}

Return a user in JSON format. For example:
```sh
curl http://localhost:8000/service/v1/users/4
```
Will return this:
```json
{
	"name": "John",
	"age": 0,
	"introduction": ""
}
```

### POST /service/v1/users/{user_id}

Modify an existing user. You do not have to send all fields, for example, to
change only the age of the user 4:

```sh
curl http://localhost:8000/service/v1/users/4 --data '{"age": 11}'
```

### DELETE /service/v1/users/{user_id}

Delete an existing user:

```sh
curl -X DELETE http://localhost:8000/service/v1/users/4
```

# Interceptors

## Interceptor Log
Log all HTTP requests to stdout in this form:

```
2016/02/20 11:09:17 GET	/favicon.ico	404	59B
2016/02/20 11:09:34 GET	/service/v1/	405	68B
2016/02/20 11:09:46 GET	/service/v1/doc	405	68B
```

## Interceptor Error
Print JSON error in this form:

```json
{
	"status_code": 404,
	"error_code": 21,
	"description_code": "User '231223' not found."
}
```

## Interceptor User
Extract and validate user from url. If the user does not exist, a 404 will be
returned.