package domain

import "time"

type WaterGoal struct {
	ID           int64
	GoalAmountML int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type WaterLog struct {
	ID          int64
	WaterGoalID int64
	AmountML    int32
	ConsumedAt  time.Time
	CreatedAt   time.Time
}
