package helpers

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"strings"
)

func LoadHtml(htmlDoc string) (*html.Node, error) {
	return htmlquery.Parse(strings.NewReader(htmlDoc))
}

func GetElements(node *html.Node, query string) []*html.Node {
	return htmlquery.Find(node, query)
}

func IsFullChain(htmlDoc string) (bool, error) {
	html, err := LoadHtml(htmlDoc)
	if err != nil {
		return false, err
	}

	scripts := GetElements(html, "/html/head/script[@src]")
	return len(scripts) == 0, nil
}
