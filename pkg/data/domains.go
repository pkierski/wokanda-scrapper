package data

import (
	"bufio"
	"bytes"
	_ "embed"
)

var (
	//go:embed domains.txt
	domainsLines []byte

	Domains = parseDomainsLines(domainsLines)
)

func parseDomainsLines(domainsLines []byte) (result []string) {
	r := bytes.NewReader(domainsLines)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return
}
