package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/javascript"
)

// Language represents a supported programming language
type Language string

const (
	Go         Language = "go"
	Python     Language = "python"
	JavaScript Language = "javascript"
)

// Node represents a node of interest in the syntax tree (function, class, etc.)
type Node struct {
	Name      string
	Kind      string
	Content   string
	StartLine int
	EndLine   int
	Signature string
}

// Parser defines the interface for code parsing
type Parser interface {
	Parse(ctx context.Context, language Language, source []byte) ([]Node, error)
}

// TreeSitterParser implements Parser using go-tree-sitter
type TreeSitterParser struct{}

func NewTreeSitterParser() *TreeSitterParser {
	return &TreeSitterParser{}
}

func (p *TreeSitterParser) Parse(ctx context.Context, language Language, source []byte) ([]Node, error) {
	var lang *sitter.Language
	switch language {
	case Go:
		lang = golang.GetLanguage()
	case Python:
		lang = python.GetLanguage()
	case JavaScript:
		lang = javascript.GetLanguage()
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(ctx, nil, source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source: %w", err)
	}

	return p.extractSymbols(tree.RootNode(), source, language), nil
}

func (p *TreeSitterParser) extractSymbols(root *sitter.Node, source []byte, lang Language) []Node {
	var symbols []Node

	// Basic symbol extraction logic (simplified for bootstrap)
	// In a real implementation, this would use Tree-Sitter queries for high precision.
	for i := 0; i < int(root.ChildCount()); i++ {
		child := root.Child(i)
		kind := child.Type()

		if isSymbolKind(kind, lang) {
			nameNode := child.ChildByFieldName("name")
			name := "anonymous"
			if nameNode != nil {
				name = nameNode.Content(source)
			}

			symbols = append(symbols, Node{
				Name:      name,
				Kind:      kind,
				Content:   child.Content(source),
				StartLine: int(child.StartPoint().Row),
				EndLine:   int(child.EndPoint().Row),
				Signature: name, // Simplified
			})
		}

		// Recursively extract symbols from children
		symbols = append(symbols, p.extractSymbols(child, source, lang)...)
	}

	return symbols
}

func isSymbolKind(kind string, lang Language) bool {
	switch lang {
	case Go:
		return kind == "function_declaration" || kind == "method_declaration" || kind == "type_declaration"
	case Python:
		return kind == "function_definition" || kind == "class_definition"
	case JavaScript:
		return kind == "function_declaration" || kind == "class_declaration" || kind == "method_definition"
	default:
		return false
	}
}
