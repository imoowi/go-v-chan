package vch

import (
	"time"
)

type ActionType string

const (
	ActionTypeIn  ActionType = `in`
	ActionTypeOut ActionType = `out`
)

type ChAction[T any] struct {
	ActionType ActionType
	ChName     string
	ChCounter  int
	Cell       T
	CellType   string
	ActionTime time.Time
}
