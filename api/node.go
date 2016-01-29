package api

type Node struct {
	Path        string
	middlewares []Handler // to be defined
	Methods     map[string]Handler
	Children    []*Node
}

type Methods struct {
	Get     Handler // to be defined
	Post    Handler
	Put     Handler
	Delete  Handler
	Options Handler
	Patch   Handler
	Else    Handler
}

func NewNode() *Node {
	return &Node{
		Path:        "",
		middlewares: []Handler{},
		Methods:     map[string]Handler{},
		Children:    []*Node{},
	}
}

func (this *Node) AddMethod(m string, h Handler) *Node {
	// TODO: m to uppercase
	this.Methods[m] = h
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
