package runtime

import "time"

type GoalStage string
const (StageObserve GoalStage = "OBSERVE"; StageSymbolLookup GoalStage = "SYMBOL_LOOKUP"; StagePatternMatch GoalStage = "PATTERN_MATCH"; StageRetrieval GoalStage = "RETRIEVAL"; StagePlan GoalStage = "PLAN"; StageAct GoalStage = "ACT"; StageLearn GoalStage = "LEARN"; StageGarden GoalStage = "GARDEN"; StageCompleted GoalStage = "COMPLETED"; StageFailed GoalStage = "FAILED")

type Goal struct { ID, Description string; CurrentStage GoalStage; Status, AssignedAgent string; CreatedAt, UpdatedAt time.Time }
