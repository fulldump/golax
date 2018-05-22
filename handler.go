package golax

// Handler is a function that implements an HTTP method for a Node, Operation,
// Interceptor.Before and Interceptor.After items.
type Handler func(c *Context)
