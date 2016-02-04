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

func (n *Node) Method(m string, h Handler) *Node {
	// TODO: m to uppercase
	M := strings.ToUpper(m)
	n.Methods[M] = h
	return n
}

func (n *Node) Middleware(m *Middleware) *Node {
	n.middlewares = append(n.middlewares, m)
	return n
}

func (n *Node) Node(p string) *Node {
	new_node := NewNode()
	new_node.Path = p

	n.Children = append(n.Children, new_node)

	return new_node
}

func (n *Node) before() {

}
