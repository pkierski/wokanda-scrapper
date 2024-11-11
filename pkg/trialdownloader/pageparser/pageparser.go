package pageparser

import (
	"slices"

	"golang.org/x/net/html"
)

func WalkNodes(node *html.Node, f func(node *html.Node)) {
	f(node)
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		WalkNodes(c, f)
	}
}

func FindNodeDown(node *html.Node, f func(node *html.Node) bool) *html.Node {
	if f(node) {
		return node
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if n := FindNodeDown(c, f); n != nil {
			return n
		}
	}
	return nil
}

func FindNodeInSiblings(node *html.Node, f func(node *html.Node) bool) *html.Node {
	for ; node != nil; node = node.NextSibling {
		if f(node) {
			return node
		}
	}
	return nil
}

func FindAttrIndex(node *html.Node, attrName string) int {
	return slices.IndexFunc(node.Attr, func(attr html.Attribute) bool {
		return attr.Key == attrName
	})
}

func FindAttrValue(node *html.Node, name string) string {
	index := FindAttrIndex(node, name)
	if index < 0 {
		return ""
	}
	return node.Attr[index].Val
}
