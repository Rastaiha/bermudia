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
	Text            string
	KnowledgeAmount int32
	RewardSource    string
}

type Answer struct {
	UserID        int32
	QuestionID    string
	Status        AnswerStatus
	RequestedHelp bool
	FileID        sql.NullString
	Filename      sql.NullString
	TextContent   sql.NullString
	Feedback      sql.NullString
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AnswerStatus int

const (
	AnswerStatusEmpty       AnswerStatus = 0
	AnswerStatusPending     AnswerStatus = 1
	AnswerStatusCorrect     AnswerStatus = 2
	AnswerStatusWrong       AnswerStatus = 3
	AnswerStatusHalfCorrect AnswerStatus = 4
)

func CheckSubmit(question BookQuestion, answer Answer) error {
	if answer.Status == AnswerStatusCorrect || answer.Status == AnswerStatusHalfCorrect {
		return ErrSubmitToCorrectAnswer
	}
	if answer.Status == AnswerStatusPending {
		return ErrSubmitToPendingAnswer
	}
	if answer.Status == AnswerStatusWrong && question.KnowledgeAmount <= 0 {
		return ErrOneTimeSubmit
	}
	return nil
}

func CheckRequestHelp(question BookQuestion, answer Answer) error {
	err := CheckSubmit(question, answer)
	if err != nil {
		return err
	}
	if question.KnowledgeAmount <= 0 {
		return Error{
			text:   "برای این سؤال نمی توانید درخواست راهنمایی کنید.",
			reason: ErrorReasonRuleViolation,
		}
	}
	return nil
}

func GetSubmissionState(question BookQuestion, answer Answer) SubmissionState {
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
	case AnswerStatusHalfCorrect:
		status = "half-correct"
	case AnswerStatusWrong:
		status = "wrong"
	}
	return SubmissionState{
		Submittable:      CheckSubmit(question, answer) == nil,
		CanRequestHelp:   CheckRequestHelp(question, answer) == nil,
		HasRequestedHelp: answer.RequestedHelp,
		Status:           status,
		Filename:         answer.Filename.String,
		Value:            answer.TextContent.String,
		Feedback:         answer.Feedback.String,
		SubmittedAt:      submittedAt,
	}
}

var (
	ErrQuestionNotFound = Error{
		text:   "question not found",
		reason: ErrorReasonResourceNotFound,
	}
	ErrAnswerNotFound = Error{
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
	ErrOneTimeSubmit = Error{
		text:   "برای این سؤال تنها یک بار می توانید پاسخ ارسال کنید.",
		reason: ErrorReasonRuleViolation,
	}
)

type KnowledgeBar struct {
	TerritoryID string `json:"territoryId"`
	Value       int32  `json:"value"`
	Total       int32  `json:"total"`
}

type CorrectionStatus int

const (
	CorrectionStatusDraft   CorrectionStatus = 0
	CorrectionStatusPending CorrectionStatus = 1
	CorrectionStatusApplied CorrectionStatus = 2
)

var CorrectionAllowedNewStatuses = []AnswerStatus{
	AnswerStatusWrong,
	AnswerStatusHalfCorrect,
	AnswerStatusCorrect,
}

type Correction struct {
	ID         string
	QuestionId string
	UserId     int32
	NewStatus  AnswerStatus
	Feedback   string
	UpdatedAt  time.Time
}
