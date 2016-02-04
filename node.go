package golax

import "strings"

type Node struct {
	Path         string
	interceptors []*Interceptor // to be defined
	Methods      map[string]Handler
	Children     []*Node
}

type Interceptor struct {
	Before Handler
	After  Handler
}

func NewNode() *Node {
	return &Node{
		Path:         "",
		interceptors: []*Interceptor{},
		Methods:      map[string]Handler{},
		Children:     []*Node{},
	}
}

func (n *Node) Method(m string, h Handler) *Node {
	// TODO: m to uppercase
	M := strings.ToUpper(m)
	n.Methods[M] = h
	return n
}

func (n *Node) Interceptor(m *Interceptor) *Node {
	n.interceptors = append(n.interceptors, m)
	return n
}

func (n *Node) Node(p string) *Node {
	new_node := NewNode()
	new_node.Path = p

	n.Children = append(n.Children, new_node)

	return new_node
}
