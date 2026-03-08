package symbols

import (
	"context"
	"fmt"

	"orion/internal/code/parser"
	sitter "github.com/smacker/go-tree-sitter"
)

type Extractor struct {
	parser *parser.CodeParser
}

func NewExtractor(p *parser.CodeParser) *Extractor {
	return &Extractor{parser: p}
}

// Extract parses source code and identifies important symbols
func (e *Extractor) Extract(ctx context.Context, filePath string, source []byte) ([]Symbol, error) {
	tree, err := e.parser.Parse(ctx, source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	var syms []Symbol
	e.traverse(tree.RootNode(), source, filePath, &syms)

	return syms, nil
}

func (e *Extractor) traverse(node *sitter.Node, source []byte, path string, syms *[]Symbol) {
	if node == nil {
		return
	}

	kind := node.Type()
	if isSymbolic(kind) {
		name := ""
		nameNode := node.ChildByFieldName("name")
		if nameNode != nil {
			name = nameNode.Content(source)
		} else if kind == "type_declaration" {
			// Find nested name
			for i := 0; i < int(node.ChildCount()); i++ {
				if node.Child(i).Type() == "type_spec" {
					if nn := node.Child(i).ChildByFieldName("name"); nn != nil {
						name = nn.Content(source)
					}
				}
			}
		}

		if name != "" {
			*syms = append(*syms, Symbol{
				Name:      name,
				Type:      kind,
				FilePath:  path,
				StartLine: int(node.StartPoint().Row),
				EndLine:   int(node.EndPoint().Row),
				Metadata:  fmt.Sprintf("node_type:%s", kind),
			})
		}
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		e.traverse(node.Child(i), source, path, syms)
	}
}

func isSymbolic(kind string) bool {
	switch kind {
	case "function_declaration", "method_declaration", "type_declaration", "struct_type", "interface_type":
		return true
	default:
		return false
	}
}
