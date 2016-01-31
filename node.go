package golax

import "strings"

type Node struct {
	Path        string
	middlewares []*Middleware // to be defined
	Methods     map[string]Handler
	Children    []*Node
}

type Middleware struct {
	Before Handler
	After  Handler
}

func NewNode() *Node {
	return &Node{
		Path:        "",
		middlewares: []*Middleware{},
		Methods:     map[string]Handler{},
		Children:    []*Node{},
	}
}

func (this *Node) AddMethod(m string, h Handler) *Node {
	// TODO: m to uppercase
	M := strings.ToUpper(m)
	this.Methods[M] = h
	return this
}

func (this *Node) AddMiddleware(m *Middleware) *Node {
	this.middlewares = append(this.middlewares, m)
	return this
}

func (this *Node) AddNode(p string) *Node {
	new_node := NewNode()
	new_node.Path = p

	this.Children = append(this.Children, new_node)

	return new_node
}

func (this *Node) before() {

}
