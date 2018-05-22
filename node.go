package golax

import "strings"

// Node represents a path part of an URL
type Node struct {
	Interceptors         []*Interceptor
	InterceptorsDeep     []*Interceptor
	Methods              map[string]Handler
	Children             []*Node
	Documentation        Doc
	DocumentationMethods map[string]Doc
	Operations           map[string]*Operation

	_path          string
	_hasOperations bool
	_isParameter   bool
	_isRegex       bool
	_isFullPath    bool
	_parameterKey  string
}

// NewNode instances and initializes a new node
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

// Method implements an HTTP method for a handler and optionally allows
// a third documentation parameter
func (n *Node) Method(m string, h Handler, d ...Doc) *Node {
	M := strings.ToUpper(m)
	n.Methods[M] = h
	if len(d) > 0 {
		n.DocumentationMethods[M] = d[0]
	}
	return n
}

// Interceptor attaches an *Interceptor to a *Node
func (n *Node) Interceptor(m *Interceptor) *Node {
	n.Interceptors = append(n.Interceptors, m)
	return n
}

// InterceptorDeep attaches an *Interceptor to a *Node but will be executed
// after all regular interceptors.
func (n *Node) InterceptorDeep(m *Interceptor) *Node {
	n.InterceptorsDeep = append(n.InterceptorsDeep, m)
	return n
}

// Doc attaches documentation to a *Node
func (n *Node) Doc(d Doc) *Node {
	n.Documentation = d

	return n
}

// Node appends a child node
func (n *Node) Node(p string) *Node {
	newNode := NewNode()
	newNode.SetPath(p)

	n.Children = append(n.Children, newNode)

	return newNode
}

// Operation appends an operation to a *Node
func (n *Node) Operation(p string) *Operation {
	n._hasOperations = true

	newOperation := NewOperation()
	newOperation.Path = p

	n.Operations[p] = newOperation

	return newOperation
}

// SetPath modifies a node path
func (n *Node) SetPath(p string) {
	n._path = p

	if "{{*}}" == p {
		n._isFullPath = true
	} else if isParameter(p) {
		n._isParameter = true
		n._parameterKey = biTrim("{", p, "}")
	} else if isRegex(p) {
		n._isRegex = true
		n._parameterKey = biTrim("(", p, ")")
	}
}

// GetPath retrieves a node path
func (n *Node) GetPath() string {
	return n._path
}

func isParameter(path string) bool {
	return '{' == path[0] && '}' == path[len(path)-1]
}

func isRegex(path string) bool {
	return '(' == path[0] && ')' == path[len(path)-1]
}

func biTrim(l, s, r string) string {
	return strings.TrimLeft(strings.TrimRight(s, r), l)
}
