package parser

import (
	"context"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type CodeParser struct {
	parser *sitter.Parser
}

func NewCodeParser() *CodeParser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &CodeParser{parser: p}
}

func (cp *CodeParser) Parse(ctx context.Context, source []byte) (*sitter.Tree, error) {
	return cp.parser.ParseCtx(ctx, nil, source)
}
