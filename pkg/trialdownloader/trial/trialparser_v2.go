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

	firstElementPred := func(node *html.Node) bool {
		return node.DataAtom == atom.Tr &&
			pageparser.FindAttrValue(node, "class") == "category table table-striped table-bordered table-hover"
	}
	secondElementPred := func(node *html.Node) bool {
		return node.DataAtom == atom.Tr &&
			pageparser.FindAttrValue(node, "class") == "row_sklad category table table-striped table-bordered table-hover"
	}

	// find first node (first case in search results)
	node := pageparser.FindNodeDown(page, firstElementPred)

	for node != nil {
		var trial Trial

		// parse first part
		parseFirstPart(node, &trial)

		node = pageparser.FindNodeInSiblings(node, secondElementPred)
		if node == nil {
			break
		}

		// parse second part
		parseSecondPart(node, &trial)
		// add entry to result
		trials = append(trials, trial)

		node = pageparser.FindNodeInSiblings(node, firstElementPred)
	}

	return
}

func parseFirstPart(node *html.Node, trial *Trial) {
	pageparser.Write(node, os.Stdout)
	fmt.Println("--------")
}

func parseSecondPart(node *html.Node, trial *Trial) {
	pageparser.Write(node, os.Stdout)
}
