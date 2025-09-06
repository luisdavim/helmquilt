package helmquilt

type Action string

const (
	ApplyAction Action = "apply"
	DiffAction  Action = "diff"
	CheckAction Action = "check"
)
