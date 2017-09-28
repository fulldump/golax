package golax

import "strings"

type Node struct {
	Interceptors         []*Interceptor
	InterceptorsDeep     []*Interceptor
	Methods              map[string]Handler
	Children             []*Node
	Documentation        Doc
	DocumentationMethods map[string]Doc
	Operations           map[string]*Operation

	_path           string
	_has_operations bool
	_is_parameter   bool
	_is_regex       bool
	_is_fullpath    bool
	_parameter_key  string
}

func NewNode() *Node {
	return &Node{
		_path:                "",
		Interceptors:         []*Interceptor{},
		InterceptorsDeep:     []*Interceptor{},
		Methods:              map[string]Handler{},
		Children:             []*Node{},
		DocumentationMethods: map[string]Doc{},
		Operations:           map[string]*Operation{},
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

func (n *Node) InterceptorDeep(m *Interceptor) *Node {
	n.InterceptorsDeep = append(n.InterceptorsDeep, m)
	return n
}

func (n *Node) Doc(d Doc) *Node {
	n.Documentation = d

	return n
}

func (n *Node) Node(p string) *Node {
	new_node := NewNode()
	new_node.SetPath(p)

	n.Children = append(n.Children, new_node)

	return new_node
}

func (n *Node) Operation(p string) *Operation {
	n._has_operations = true

	new_operation := NewOperation()
	new_operation.Path = p

	n.Operations[p] = new_operation

	return new_operation
}

func (n *Node) SetPath(p string) {
	n._path = p

	if "{{*}}" == p {
		n._is_fullpath = true
	} else if is_parameter(p) {
		n._is_parameter = true
		n._parameter_key = bi_trim("{", p, "}")
	} else if is_regex(p) {
		n._is_regex = true
		n._parameter_key = bi_trim("(", p, ")")
	}
}

func (n *Node) GetPath() string {
	return n._path
}

func is_parameter(path string) bool {
	return '{' == path[0] && '}' == path[len(path)-1]
}

func is_regex(path string) bool {
	return '(' == path[0] && ')' == path[len(path)-1]
}

func bi_trim(l, s, r string) string {
	return strings.TrimLeft(strings.TrimRight(s, r), l)
}
