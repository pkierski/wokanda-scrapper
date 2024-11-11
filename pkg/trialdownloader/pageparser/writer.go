package pageparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

func Write(node *html.Node, w io.Writer) {
	writeNode(node, w, 0)
}

func writeNode(node *html.Node, w io.Writer, indent int) {
	if true || node.DataAtom != 0 || strings.TrimSpace(node.Data) != "" {
		fmt.Fprintf(w, "%vtype: '%v', data: '%v', attr: %v\n",
			strings.Repeat("  ", indent), node.DataAtom, node.Data, node.Attr)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		writeNode(c, w, indent+1)
	}
}
