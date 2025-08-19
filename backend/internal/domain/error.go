package domain

const (
	ErrorReasonResourceNotFound = iota
	ErrorReasonRuleViolation
)

type Error struct {
	text   string
	reason int
}

func (e Error) Error() string {
	return e.text
}

func (e Error) Reason() int {
	return e.reason
}
