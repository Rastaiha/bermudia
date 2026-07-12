package domain

const (
	ErrorReasonResourceNotFound = iota
	ErrorReasonRuleViolation
)

type Error struct {
	text   string
	reason int
}

// NewError builds a domain.Error with the given reason and message. It exists so
// packages outside domain (e.g. service) can return typed domain errors.
func NewError(reason int, text string) Error {
	return Error{reason: reason, text: text}
}

func (e Error) Error() string {
	return e.text
}

func (e Error) Reason() int {
	return e.reason
}
