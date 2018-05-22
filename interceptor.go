package golax

// Interceptor are pieces of code attached to nodes that can interact with the
// context and break the execution.
// A interceptor has two pieces of code, `Before` is executed before the handler
// and `After` is executed after the handler. The `After` code is executed
// always if the `Before` code has been executed successfully (without calling
// to `context.Error(...)`.
type Interceptor struct {
	Before        Handler
	After         Handler
	Documentation Doc
}
