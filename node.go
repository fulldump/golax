package golax

import "strings"

type Node struct {
	Path                 string
	Interceptors         []*Interceptor
	Methods              map[string]Handler
	Children             []*Node
	Documentation        Doc
	DocumentationMethods map[string]Doc
}

func NewNode() *Node {
	return &Node{
		Path:                 "",
		Interceptors:         []*Interceptor{},
		Methods:              map[string]Handler{},
		Children:             []*Node{},
		DocumentationMethods: map[string]Doc{},
	}
}

func (n *Node) Method(m string, h Handler, d ...Doc) *Node {
	M := strings.ToUpper(m)
	n.Methods[M] = h
	if len(d) > 0 {
		n.DocumentationMethods[M] = d[0]
	}
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
