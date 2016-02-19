package golax

import "strings"

type Node struct {
	Path          string
	Interceptors  []*Interceptor
	Methods       map[string]Handler
	Children      []*Node
	Documentation Doc
}

func NewNode() *Node {
	return &Node{
		Path:         "",
		Interceptors: []*Interceptor{},
		Methods:      map[string]Handler{},
		Children:     []*Node{},
	}
}

func (n *Node) Method(m string, h Handler) *Node {
	M := strings.ToUpper(m)
	n.Methods[M] = h
	return n
}

func (n *Node) Interceptor(m *Interceptor) *Node {
	n.Interceptors = append(n.Interceptors, m)
	return n
}

func (n *Node) Doc(d Doc) *Node {
	n.Documentation = d

	return n
}

func (n *Node) Node(p string) *Node {
	new_node := NewNode()
	new_node.Path = p

	n.Children = append(n.Children, new_node)

	return new_node
}
