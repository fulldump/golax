package golax

// Doc represents documentation information that can be attached to a node,
// method or interceptor.
type Doc struct {
	Name        string
	Description string
	Ommit       bool
}
