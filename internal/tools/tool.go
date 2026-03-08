package tools

import "context"

type Tool interface { Name() string; Description() string; Execute(ctx context.Context, input interface{}) (interface{}, error) }

type BaseTool struct { name, desc string }
func (t *BaseTool) Name() string { return t.name }
func (t *BaseTool) Description() string { return t.desc }
