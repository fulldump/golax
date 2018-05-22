package golax

import (
	"encoding/json"
	"log"
)

// InterceptorError prints an error in JSON format if Context.LastError is
// not nil.
// Example:
//     {
//         "status_code": 404,
//         "error_code": 1000023,
//         "description_code": "User 'fulanez' not found.",
//     }
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

// InterceptorLog prints an access log to standard output.
// Example:
// 2016/02/20 11:09:17 GET	/favicon.ico	404	59B
// 2016/02/20 11:09:34 GET	/service/v1/	405	68B
// 2016/02/20 11:09:46 GET	/service/v1/doc	405	68B
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

// InterceptorNoCache set some headers to force response to not be cached
// by user agent.
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
