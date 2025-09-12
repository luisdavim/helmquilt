package helmquilt

import "errors"

type Action string

const (
	ApplyAction Action = "apply"
	DiffAction  Action = "diff"
	CheckAction Action = "check"
)

var (
	ErrChartsChanged = errors.New("not all charts are up to date")
	ErrDirtyCharts   = errors.New("some charts need to be cleaned")
)
