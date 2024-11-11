package trial

import (
	"bytes"
	"fmt"
	"os"

	"github.com/pkierski/wokanda-scrapper/pkg/trialdownloader/pageparser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// ParseV1 parses one page from type pages like
// "https://bialystok.sa.gov.pl/zalatw-sprawe/e-wokanda".
func ParseV2(data []byte) (trials []Trial, err error) {
	page, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("parsing trial page: %w", err)
		return
	}

	// find first node (first case in search results)

	node := pageparser.FindNodeDown(page, func(node *html.Node) bool {
		return node.DataAtom == atom.Tr &&
			pageparser.FindAttrValue(node, "class") == "category table table-striped table-bordered table-hover"
	})

	pageparser.Write(node, os.Stdout)
	fmt.Println("--------")

	node = pageparser.FindNodeInSiblings(node, func(node *html.Node) bool {
		return node.DataAtom == atom.Tr &&
			pageparser.FindAttrValue(node, "class") == "row_sklad category table table-striped table-bordered table-hover"
	})

	pageparser.Write(node, os.Stdout)

	return nil, nil
}
