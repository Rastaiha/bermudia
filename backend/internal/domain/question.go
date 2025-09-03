package domain

import (
	"database/sql"
	"time"
)

type Question struct {
	ID          string   `json:"id"`
	Text        string   `json:"text"`
	InputType   string   `json:"inputType"`
	InputAccept []string `json:"inputAccept"`
}

type BookQuestion struct {
	QuestionID      string
	BookID          string
	KnowledgeAmount int32
	RewardSource    string
}

type Answer struct {
	UserID      int32
	QuestionID  string
	Status      AnswerStatus
	FileID      sql.NullString
	Filename    sql.NullString
	TextContent sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AnswerStatus int

const (
	AnswerStatusEmpty   AnswerStatus = 0
	AnswerStatusPending AnswerStatus = 1
	AnswerStatusCorrect AnswerStatus = 2
	AnswerStatusWrong   AnswerStatus = 3
)

func GetSubmissionStateFromAnswer(answer Answer) SubmissionState {
	submittedAt := answer.UpdatedAt.UnixMilli()
	if answer.Status == AnswerStatusEmpty {
		submittedAt = 0
	}
	status := ""
	switch answer.Status {
	case AnswerStatusEmpty:
		status = "empty"
	case AnswerStatusPending:
		status = "pending"
	case AnswerStatusCorrect:
		status = "correct"
	case AnswerStatusWrong:
		status = "wrong"
	}
	return SubmissionState{
		Submittable: answer.Status == AnswerStatusEmpty || answer.Status == AnswerStatusWrong,
		Status:      status,
		Filename:    answer.Filename.String,
		Value:       answer.TextContent.String,
		SubmittedAt: submittedAt,
	}
}

var (
	ErrQuestionNotFound = Error{
		text:   "question not found",
		reason: ErrorReasonResourceNotFound,
	}
	ErrSubmitToNonExistingAnswer = Error{
		text:   "answer id does not exist",
		reason: ErrorReasonResourceNotFound,
	}
	ErrSubmitToCorrectAnswer = Error{
		text:   "به این سؤال یک بار پاسخ داده اید.",
		reason: ErrorReasonRuleViolation,
	}
	ErrSubmitToPendingAnswer = Error{
		text:   "در حال حاضر یک پاسخ بررسی نشده برای این سؤال وجود دارد.",
		reason: ErrorReasonRuleViolation,
	}
)

type KnowledgeBar struct {
	TerritoryID string `json:"territoryId"`
	Value       int32  `json:"value"`
	Total       int32  `json:"total"`
}

type Correction struct {
	ID         string
	QuestionId string
	UserId     int32
	IsCorrect  bool
	CreatedAt  time.Time
}
