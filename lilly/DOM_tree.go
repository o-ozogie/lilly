package lilly

import (
	"golang.org/x/net/html"
)

type DOMTree struct {
	RootNode *Node

	ContentNode   *Node
	ThresholdNode *Node
}

func NewDOMTree(root *html.Node) *DOMTree {
	_DOMTree := &DOMTree{}
	_DOMTree.RootNode = _DOMTree.NewRootNode(root)
	_DOMTree.RootNode.setCompositeTextDensity()
	_DOMTree.setThresholdNode()

	return _DOMTree
}

func (_DOMTree *DOMTree) NewRootNode(root *html.Node) *Node {
	rootNode := &Node{Node: root, DOMTree: _DOMTree}

	rootNode.setDOMTree()
	rootNode.setTagCount()
	rootNode.setCharacterCount()
	rootNode.setTextDensity()

	return rootNode
}

func (_DOMTree *DOMTree) setThresholdNode() {
	_DOMTree.RootNode.Traversal(func(node *Node) {
		if _DOMTree.ContentNode == nil || _DOMTree.ContentNode.CompositeDensitySum < node.CompositeDensitySum {
			_DOMTree.ContentNode = node
		}
	})
	for parent := _DOMTree.ContentNode.Parent; parent != nil; parent = parent.Parent {
		if _DOMTree.ThresholdNode == nil || _DOMTree.ThresholdNode.CompositeTextDensity > parent.CompositeTextDensity {
			_DOMTree.ThresholdNode = parent
		}
	}
}

func (_DOMTree *DOMTree) ExtractContent() (content string) {
	_DOMTree.RootNode.Traversal(func(node *Node) {
		if node.Type == html.ElementNode && node.CompositeTextDensity >= _DOMTree.ThresholdNode.CompositeTextDensity {
			for _, child := range node.Children {
				if child.Type == html.TextNode {
					content += child.Data
				}
			}
		}
	})
	return
}

func (_DOMTree *DOMTree) ExtractAccurateContent() (content string) {
	_DOMTree.ContentNode.Traversal(func(node *Node) {
		if node.Type == html.TextNode {
			content += node.Data
		}
	})
	return
}
