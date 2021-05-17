package lilly

import (
	"bytes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
)

func NewDOMTreeFromURL(url string) *DOMTree {
	response, e := http.Get(url)
	if e != nil {
		panic(e)
	}
	if response.StatusCode/100 != 2 {
		panic("response not valid")
	}
	defer response.Body.Close()

	var body bytes.Buffer
	converted, e := charset.NewReaderLabel(getLabel(io.TeeReader(response.Body, &body)), &body)
	if e != nil {
		panic(e)
	}

	document, e := html.Parse(converted)
	if e != nil {
		panic(e)
	}

	return NewDOMTree(getNode(document, atom.Body))
}

func getLabel(body io.Reader) string {
	label := "utf-8"
	document, e := html.Parse(body)
	if e != nil {
		panic(e)
	}

	for child := getNode(document, atom.Head).FirstChild; child != nil; child = child.NextSibling {
		if child.DataAtom == atom.Meta {
			for _, attribute := range child.Attr {
				if attribute.Key == "charset" {
					label = attribute.Val
				}
			}
		}
	}

	return label
}

func getNode(document *html.Node, dataAtom atom.Atom) (found *html.Node) {
	var getNode func(node *html.Node)
	getNode = func(node *html.Node) {
		if node.DataAtom == dataAtom {
			found = node
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			getNode(child)
		}
	}
	getNode(document)

	if found != nil {
		return
	}

	panic("cannot find <" + dataAtom.String() + ">")
}
