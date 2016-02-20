package apidoc

import (
	"fmt"
	"strings"

	"github.com/fulldump/golax"
)

func Build(a *golax.Api) {

	doc := a.Root.
		Node("doc").
		Doc(golax.Doc{Ommit: true, Description: `
This API sub tree is dedicated to self-documentation. It contains documentation
in several formats: markdown, html and swagger at this moment.

**Permissions**:
Anyone can read this.
		`})

	md := doc.
		Node("md").
		Doc(golax.Doc{Description: `
Documentation in markdown format.
		`})
	md.Method("GET", func(c *golax.Context) {
		PrintApiMd(NodePrint{
			Api:             a,
			Node:            a.Root,
			Context:         c,
			Path:            a.Prefix,
			AllInterceptors: map[*golax.Interceptor]*golax.Interceptor{},
		})
	})

	html := doc.
		Node("html").
		Doc(golax.Doc{Description: `
Documentation in html format.
		`})
	html.Method("GET", func(c *golax.Context) {
		c.Error(501, "Unimplemented")
	})

	swagger := doc.
		Node("swagger").
		Doc(golax.Doc{Description: `
Documentation in swagger format.
		`})
	swagger.Method("GET", func(c *golax.Context) {
		c.Error(501, "Unimplemented")
	})

}

type NodePrint struct {
	Api                 *golax.Api
	Node                *golax.Node
	Context             *golax.Context
	Path                string
	Level               int
	AllInterceptors     map[*golax.Interceptor]*golax.Interceptor
	CurrentInterceptors []*golax.Interceptor
}

func (np *NodePrint) Println(s string) {
	fmt.Fprintln(np.Context.Response, s)
}

func (np *NodePrint) Printf(f string, a ...interface{}) {
	fmt.Fprintf(np.Context.Response, f, a...)
}

func md_link(l string) string {
	l = strings.ToLower(l)
	l = strings.Replace(l, " ", "-", -1)
	l = strings.Replace(l, "/", "", -1)
	l = strings.Replace(l, "{", "", -1)
	l = strings.Replace(l, "}", "", -1)
	return "#" + l
}

func md_description(d string) string {
	d = strings.TrimSpace(d)
	d = strings.Replace(d, "\n´´´", "\n```", -1)
	return d
}

func PrintApiMd(p NodePrint) {

	if p.Node.Documentation.Ommit {
		return
	}

	for _, i := range p.Node.Interceptors {
		p.AllInterceptors[i] = i
		p.CurrentInterceptors = append(p.CurrentInterceptors, i)
	}

	is_root := 0 == p.Level
	p.Level++

	// Title
	if is_root {
		p.Println("# API Documentation")
	} else {
		p.Path += "/" + p.Node.Path
		p.Println("\n## " + p.Path + "\n")
	}

	// Description
	p.Println(md_description(p.Node.Documentation.Description))

	// Applied interceptors
	if len(p.CurrentInterceptors) > 0 {
		interceptors := "\n**Interceptors chain:** "
		if is_root {
			interceptors = "\n**Interceptors applied to all API:** "
		}
		for _, v := range p.CurrentInterceptors {
			name := v.Documentation.Name
			link := md_link("Interceptor " + name)
			interceptors += " [`" + name + "`](" + link + ") "
		}
		p.Println(interceptors)
	}

	// Implemented methods
	if len(p.Node.Methods) > 0 {
		methods := "\n**Methods:** "
		for k, _ := range p.Node.Methods {
			link := md_link(k + " " + p.Path)
			methods += " [`" + k + "`](" + link + ") "
		}
		p.Println(methods)

		for k, _ := range p.Node.Methods {
			methods += " `" + k + "` "
			p.Println("\n### " + k + " " + p.Path + "\n")

			if d, e := p.Node.DocumentationMethods[k]; e {
				p.Println(md_description(d.Description))
			}
		}
	}

	// Document children
	for _, child := range p.Node.Children {
		p.Node = child
		PrintApiMd(p)
	}

	if is_root {
		p.Println("\n# Interceptors")

		for _, v := range p.AllInterceptors {
			p.Println("\n## Interceptor " + v.Documentation.Name)
			p.Println(md_description(v.Documentation.Description))
		}
	}
}
