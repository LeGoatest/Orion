package cognition

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"orion/internal/agents"
	"orion/internal/types"
	"orion/internal/workspace"
)

type CognitionEngine struct {
	pipeline     OODALPipeline
	eventBus     *types.EventBus
	dispatcher   *agents.Dispatcher
	workspaceMgr *workspace.Manager
}

func NewCognitionEngine(p OODALPipeline, eb *types.EventBus, ad *agents.Dispatcher, wm *workspace.Manager) *CognitionEngine {
	ce := &CognitionEngine{
		pipeline:     p,
		eventBus:     eb,
		dispatcher:   ad,
		workspaceMgr: wm,
	}

	// Subscribe to completion events to drive the multi-agent pipeline
	go ce.listenForEvents()

	return ce
}

func (ce *CognitionEngine) listenForEvents() {
	ch := ce.eventBus.Subscribe("*")
	for event := range ch {
		ce.HandleStageCompletion(context.Background(), event)
	}
}

func (ce *CognitionEngine) Process(ctx context.Context, workspaceID string, goalID string, description string) error {
	fmt.Printf("CognitionEngine: Initializing goal %s in workspace %s\n", goalID, workspaceID)

	db, err := ce.workspaceMgr.GetWorkspaceDB(workspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace db: %w", err)
	}

	// Initial Goal Creation
	now := time.Now()
	_, err = db.Exec(`INSERT INTO goals (id, description, current_stage, status, created_at, updated_at)
					  VALUES (?, ?, ?, ?, ?, ?)`,
		goalID, description, "OBSERVE", "active", now, now)
	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}

	// First stage: OBSERVE
	return ce.dispatcher.Dispatch(ctx, db, workspaceID, "OBSERVE", goalID, description)
}

func (ce *CognitionEngine) HandleStageCompletion(ctx context.Context, event types.Event) {
	// Only handle stage completion events
	payload, ok := event.Payload.(map[string]interface{})
	if !ok {
		return
	}
	goalID := event.GoalID
	workspaceID := event.WorkspaceID

	if goalID == "" || workspaceID == "" {
		return
	}

	var db *sql.DB
	var err error
	db, err = ce.workspaceMgr.GetWorkspaceDB(workspaceID)
	if err != nil {
		fmt.Printf("Cognition: Failed to get db for workspace %s: %v\n", workspaceID, err)
		return
	}

	switch event.Type {
	case "cognition.OBSERVE.completed":
		fmt.Printf("Pipeline [%s]: OBSERVE complete -> SYMBOL_LOOKUP\n", goalID)
		ce.dispatcher.Dispatch(ctx, db, workspaceID, "SYMBOL_LOOKUP", goalID, payload["result"])

	case "cognition.SYMBOL_LOOKUP.completed":
		fmt.Printf("Pipeline [%s]: SYMBOL_LOOKUP complete -> PATTERN_MATCH\n", goalID)
		ce.dispatcher.Dispatch(ctx, db, workspaceID, "PATTERN_MATCH", goalID, payload["result"])

	case "cognition.PATTERN_MATCH.completed":
		fmt.Printf("Pipeline [%s]: PATTERN_MATCH complete -> RETRIEVAL\n", goalID)
		ce.dispatcher.Dispatch(ctx, db, workspaceID, "RETRIEVAL", goalID, payload["result"])

	case "cognition.RETRIEVAL.completed":
		fmt.Printf("Pipeline [%s]: RETRIEVAL complete -> PLAN\n", goalID)
		ce.dispatcher.Dispatch(ctx, db, workspaceID, "PLAN", goalID, payload["result"])

	case "cognition.PLAN.completed":
		fmt.Printf("Pipeline [%s]: PLAN complete -> ACT\n", goalID)
		ce.dispatcher.Dispatch(ctx, db, workspaceID, "ACT", goalID, payload["result"])

	case "cognition.ACT.completed":
		fmt.Printf("Pipeline [%s]: ACT complete -> LEARN\n", goalID)
		ce.eventBus.Publish(types.Event{
			Type:        "cognition.LEARN.completed",
			GoalID:      goalID,
			WorkspaceID: workspaceID,
			Payload:     map[string]interface{}{"result": payload["result"]},
			CreatedAt:   time.Now(),
		})

	case "cognition.LEARN.completed":
		fmt.Printf("Pipeline [%s]: LEARN complete -> GARDEN (scheduled)\n", goalID)
		if db != nil {
			_, _ = db.Exec("UPDATE goals SET current_stage = ?, status = ?, updated_at = ? WHERE id = ?", "COMPLETED", "completed", time.Now(), goalID)
		}
		ce.eventBus.Publish(types.Event{
			Type:        "cognition.GARDEN.scheduled",
			GoalID:      goalID,
			WorkspaceID: workspaceID,
			CreatedAt:   time.Now(),
		})
		ce.eventBus.Publish(types.Event{
			Type:        "cognition.goal.completed",
			GoalID:      goalID,
			WorkspaceID: workspaceID,
			CreatedAt:   time.Now(),
		})
	}
}

type OODALPipeline interface {
	Observe(context.Context, interface{}) (interface{}, error)
	Orient(context.Context, interface{}) (interface{}, error)
	Decide(context.Context, interface{}) (interface{}, error)
	Act(context.Context, interface{}) (interface{}, error)
	Learn(context.Context, interface{}) error
	Garden(context.Context, string) error
}

type DefaultPipeline struct{}

func (p *DefaultPipeline) Observe(ctx context.Context, i interface{}) (interface{}, error) { return i, nil }
func (p *DefaultPipeline) Orient(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Decide(ctx context.Context, i interface{}) (interface{}, error)  { return i, nil }
func (p *DefaultPipeline) Act(ctx context.Context, i interface{}) (interface{}, error)     { return i, nil }
func (p *DefaultPipeline) Learn(ctx context.Context, i interface{}) error                 { return nil }
func (p *DefaultPipeline) Garden(ctx context.Context, s string) error                     { return nil }
