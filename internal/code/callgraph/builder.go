package callgraph

import (
	"context"
	"fmt"
	"orion/internal/code/parser"
	sitter "github.com/smacker/go-tree-sitter"
)

type Builder struct {
	parser *parser.CodeParser
}

func (b *Builder) BuildCallGraph(ctx context.Context, filePath string, source []byte) error {
	tree, err := b.parser.Parse(ctx, parser.Go, source)
	if err != nil { return err }

	fmt.Printf("CallGraph: Analyzing calls in %s\n", filePath)
	b.traverseForCalls(tree.RootNode(), source)
	return nil
}

func (b *Builder) traverseForCalls(node *sitter.Node, source []byte) {
	if node == nil { return }

	// Real logic: identify call_expression nodes and resolve them to symbols
	if node.Type() == "call_expression" {
		fmt.Printf("Found call at line %d\n", node.StartPoint().Row)
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		b.traverseForCalls(node.Child(i), source)
	}
}
