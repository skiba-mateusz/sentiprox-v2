package detector

type Action string

const (
	ActionPass  Action = "pass"
	ActionBlock Action = "block"
	ActionFlag  Action = "flag"
)

type Source string

const (
	SourceNone       Source = ""
	SourceBehavioral Source = "behavioral"
	SourceML         Source = "ml"
	SourceSystem     Source = "sytem"
)

const (
	CategoryBenign  = "benign"
	CategoryUnknown = "unknown"
)

type Decision struct {
	Action   Action
	Source   Source
	Category string
	Status   int
}

var pass = Decision{Action: ActionPass}

func block(src Source, cat string, status int) Decision {
	return Decision{Action: ActionBlock, Source: src, Category: cat, Status: status}
}

func flag(src Source, cat string) Decision {
	return Decision{Action: ActionFlag, Source: src, Category: cat}
}

func (d Decision) Blocked() bool {
	return d.Action == ActionBlock
}
