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
	if len(scripts) > 0 {
		return false, nil
	}

	scripts = GetElements(html, "/html/body/script[@src]")
	if len(scripts) > 0 {
		return false, nil
	}

	links := GetElements(html, "/html/head/link[@href]")
	if len(links) > 0 {
		return false, nil
	}
	links = GetElements(html, "/html/body/link[@href]")
	if len(links) > 0 {
		return false, nil
	}
	return true, nil
}
