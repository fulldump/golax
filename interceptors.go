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
	Documentation: Doc{
		Name: "Error",
		Description: `
Print JSON error in this form:

´´´json
{
	"status_code": 404,
	"error_code": 21,
	"description_code": "User '231223' not found."
}
´´´
		`,
	},
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
	Documentation: Doc{
		Name: "Log",
		Description: `
Log all HTTP requests to stdout in this form:

´´´
2016/02/20 11:09:17 GET	/favicon.ico	404	59B
2016/02/20 11:09:34 GET	/service/v1/	405	68B
2016/02/20 11:09:46 GET	/service/v1/doc	405	68B
´´´
		`,
	},
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

/**
 * `InterceptorNoCache`
 * Send headers to disable browser caching
 */
var InterceptorNoCache = &Interceptor{
	Documentation: Doc{
		Name: "InterceptorNoCache",
		Description: `
			Avoid caching via http headers
		`,
	},
	Before: func(c *Context) {
		add := c.Response.Header().Add

		add("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		add("Pragma", "no-cache")
		add("Expires", "0")
	},
}
