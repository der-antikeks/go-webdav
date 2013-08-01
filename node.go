package webdav

import (
	"encoding/xml"
	"io"
)

type Node struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []*Node
	Parent   *Node
}

func NodeFromXml(r io.Reader) (*Node, error) {
	var cur, parent *Node

	decoder := xml.NewDecoder(r)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		switch tok := token.(type) {
		case xml.StartElement:
			parent = cur

			// if tok.Name.Space != "DAV:" {
			if tok.Name.Space == "" {
				return nil, ErrMalformedXml
			}

			cur = &Node{
				Name:   tok.Name,
				Attr:   tok.Attr,
				Parent: parent,
			}

			if parent != nil {
				parent.Children = append(parent.Children, cur)
			}
		case xml.EndElement:
			if cur.Parent == nil {
				return cur, nil
			}
			cur = cur.Parent
		default:
			//log.Printf("%T", tok)
		}
	}

	return cur, nil
}

func (n Node) HasChildren(name string) bool {
	for _, v := range n.Children {
		if v.Name.Local == name {
			return true
		}
	}
	return false
}

func (n *Node) GetChildrens(name string) []*Node {
	if name == "*" {
		return n.Children
	}

	var ret []*Node

	for _, v := range n.Children {
		if v.Name.Local == name {
			ret = append(ret, v)
		}
	}

	return ret
}

func (n *Node) FirstChildren(name string) *Node {
	if name == "*" && len(n.Children) > 0 {
		return n.Children[0]
	}

	for _, v := range n.Children {
		if v.Name.Local == name {
			return v
		}
	}

	return nil
}

func (n *Node) String() string {
	r := "<" + n.Name.Local + " xmlns=\"" + n.Name.Space + "\">\n"

	for _, v := range n.Children {
		r += v.String() + "\n"
	}

	r += "<\\" + n.Name.Local + ">"

	return r
}
