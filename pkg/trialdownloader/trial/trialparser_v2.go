package trial

import (
	"bytes"
	"fmt"

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
	var node *html.Node
	pageparser.WalkNodes(page, func(n *html.Node) {
		if node != nil {
			return
		}
		if n.DataAtom == atom.Tr {
			if pageparser.FindAttrValue(n, "class") == "category table table-striped table-bordered table-hover" {
				node = n
			}
		}

	})

	return nil, nil
}
