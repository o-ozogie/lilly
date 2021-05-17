package lilly

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"math"
	"strings"
	"unicode/utf8"
)

type Node struct {
	*html.Node
	DOMTree *DOMTree

	Parent   *Node
	Children []*Node

	TagCount             float64
	LinkTagCount         float64
	CharacterCount       float64
	LinkCharacterCount   float64
	TextDensity          float64
	CompositeTextDensity float64
	CompositeDensitySum  float64
}

type Action func(node *Node)

func (node *Node) Traversal(action Action) {
	action(node)
	for _, child := range node.Children {
		child.Traversal(action)
	}
}

func (node *Node) setDOMTree() {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasInvalidType(child) || hasInvalidData(child) || hasInvalidDataAtom(child) {
			continue
		}
		childNode := &Node{
			DOMTree: node.DOMTree,
			Node:    child,
			Parent:  node,
		}
		childNode.setDOMTree()
		node.Children = append(node.Children, childNode)
	}
}

func hasInvalidType(node *html.Node) bool {
	return node.Type != html.ElementNode && node.Type != html.TextNode
}

func hasInvalidData(node *html.Node) bool {
	var filtered string

	filtered = strings.ReplaceAll(node.Data, "\n", "")
	filtered = strings.ReplaceAll(filtered, " ", "")
	filtered = strings.ReplaceAll(filtered, "\t", "")

	return len(filtered) == 0
}

func hasInvalidDataAtom(node *html.Node) bool {
	return node.DataAtom == atom.Style || node.DataAtom == atom.Script
}

func (node *Node) setTagCount() {
	for _, child := range node.Children {
		if child.Type != html.ElementNode {
			continue
		}

		child.setTagCount()

		node.TagCount += child.TagCount
		node.TagCount++
		node.LinkTagCount += child.LinkTagCount

		if child.TagCount == 0 {
			child.TagCount = 1
		}
	}
	if node.DataAtom == atom.A {
		node.LinkTagCount++
	}
}

func (node *Node) setCharacterCount() {
	for _, child := range node.Children {
		child.setCharacterCount()

		node.CharacterCount += child.CharacterCount
		node.LinkCharacterCount += child.LinkCharacterCount
	}
	if node.Type == html.TextNode {
		node.CharacterCount += float64(utf8.RuneCountInString(node.Data))
		if node.Parent.DataAtom == atom.A {
			node.LinkCharacterCount += float64(utf8.RuneCountInString(node.Data))
		}
	}
}

func (node *Node) setTextDensity() {
	for _, child := range node.Children {
		if child.Type == html.ElementNode {
			child.setTextDensity()
		}
	}
	node.TextDensity = node.CharacterCount / node.TagCount
}

func (node *Node) setCompositeTextDensity() {
	for _, child := range node.Children {
		if child.Type == html.ElementNode {
			child.setCompositeTextDensity()

			node.CompositeDensitySum += child.CompositeTextDensity
		}
	}
	node.CompositeTextDensity = func() float64 {
		linkCharacterCount := node.LinkCharacterCount
		if linkCharacterCount == 0 {
			linkCharacterCount = 1
		}
		linkTagCount := node.LinkTagCount
		if linkTagCount == 0 {
			linkTagCount = 1
		}
		notLinkCharacterCount := node.CharacterCount - node.LinkCharacterCount
		if notLinkCharacterCount == 0 {
			notLinkCharacterCount = 1
		}
		bodyCharacterCount := node.DOMTree.RootNode.CharacterCount
		if bodyCharacterCount == 0 {
			bodyCharacterCount = 1
		}

		hyperlinkRate := (node.CharacterCount / linkCharacterCount) * (node.TagCount / linkTagCount)
		textRate := node.CharacterCount/notLinkCharacterCount*node.LinkCharacterCount +
			node.DOMTree.RootNode.LinkCharacterCount/bodyCharacterCount*node.CharacterCount + math.E

		compositeTextDensity := node.TextDensity * (math.Log(hyperlinkRate) / math.Log(math.Log(textRate)))
		if math.IsNaN(compositeTextDensity) {
			return 0
		}
		return compositeTextDensity
	}()
}
