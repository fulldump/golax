package api

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request, c *Context)
