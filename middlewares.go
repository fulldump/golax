package golax

/**
 * Core middlewares free to use
 */

var MiddlewareError = &Middleware{
	After: func(c *Context) {
		if nil != c.LastError {
			c.Response.WriteHeader(c.LastError.Status)
		}
	},
}
