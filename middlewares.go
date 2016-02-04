package golax

import "encoding/json"

/**
 * Core middlewares free to use
 */

var InterceptorError = &Interceptor{
	After: func(c *Context) {
		if nil != c.LastError {
			json.NewEncoder(c.Response).Encode(c.LastError)
		}
	},
}
