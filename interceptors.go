package golax

import (
	"encoding/json"
	"log"
)

/**
 * Core interceptors free to use
 */

/**
 * `InterceptorError`
 * Print a JSON with the last error if exists
 */
var InterceptorError = &Interceptor{
	After: func(c *Context) {
		if nil != c.LastError {
			json.NewEncoder(c.Response).Encode(c.LastError)
		}
	},
}

/**
 * `InterceptorLog`
 * Log request and response
 */
var InterceptorLog = &Interceptor{
	After: func(c *Context) {
		log.Printf(
			"%s\t%s\t%d\t%dB",
			c.Request.Method,
			c.Request.URL.RequestURI(),
			c.Response.StatusCode,
			c.Response.Length,
		)
	},
}
