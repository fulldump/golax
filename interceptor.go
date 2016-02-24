package golax

type Interceptor struct {
	Before        Handler
	After         Handler
	Documentation Doc
}
