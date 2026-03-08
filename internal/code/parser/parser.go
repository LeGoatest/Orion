package parser

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/javascript"
)

type Language string

const (
	Go         Language = "go"
	Python     Language = "python"
	TypeScript Language = "typescript"
)

type CodeParser struct {
	parsers map[Language]*sitter.Parser
}

func NewCodeParser() *CodeParser {
	cp := &CodeParser{
		parsers: make(map[Language]*sitter.Parser),
	}

	pGo := sitter.NewParser()
	pGo.SetLanguage(golang.GetLanguage())
	cp.parsers[Go] = pGo

	pPy := sitter.NewParser()
	pPy.SetLanguage(python.GetLanguage())
	cp.parsers[Python] = pPy

	pTs := sitter.NewParser()
	// Using javascript language as a close proxy for typescript in this bootstrap phase
	pTs.SetLanguage(javascript.GetLanguage())
	cp.parsers[TypeScript] = pTs

	return cp
}

func (cp *CodeParser) Parse(ctx context.Context, lang Language, source []byte) (*sitter.Tree, error) {
	parser, ok := cp.parsers[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
	return parser.ParseCtx(ctx, nil, source)
}
