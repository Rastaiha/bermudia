package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
	"slices"
	"strings"
	"time"
)

type Correction struct {
	cfg           config.Config
	questionStore domain.QuestionStore
}

func NewCorrection(cfg config.Config, questionStore domain.QuestionStore) *Correction {
	return &Correction{
		cfg:           cfg,
		questionStore: questionStore,
	}
}

func (c *Correction) AutoCorrect(ctx context.Context, answer domain.Answer) bool {
	correction := domain.Correction{
		ID:         domain.NewID(domain.ResourceTypeCorrection),
		QuestionId: answer.QuestionID,
		UserId:     answer.UserID,
		UpdatedAt:  time.Now().UTC(),
	}
	create := false
	lowerFilename := strings.ToLower(answer.Filename.String)
	if strings.Contains(lowerFilename, "false") || strings.Contains(answer.TextContent.String, "false") || strings.Contains(answer.TextContent.String, "0") {
		correction.NewStatus = domain.AnswerStatusWrong
		correction.Feedback = "یه کم بیشتر فکر کن!"
		create = true
	} else if strings.Contains(lowerFilename, "true") || strings.Contains(answer.TextContent.String, "true") || strings.Contains(answer.TextContent.String, "1") {
		correction.NewStatus = domain.AnswerStatusCorrect
		correction.Feedback = "آفرین، حالا بریم سؤال بعدی :)"
		create = true
	} else if strings.Contains(lowerFilename, "half") || strings.Contains(answer.TextContent.String, "half") || strings.Contains(answer.TextContent.String, "2") {
		correction.NewStatus = domain.AnswerStatusHalfCorrect
		correction.Feedback = "تقریباً درسته :)"
		create = true
	}
	if create {
		if err := c.questionStore.CreateCorrection(ctx, correction); err != nil {
			slog.Error("failed to create correction", slog.String("error", err.Error()))
		} else if err := c.questionStore.FinalizeCorrection(ctx, correction.ID); err != nil {
			slog.Error("failed to finalize correction", slog.String("error", err.Error()))
		}
	}
	return create
}

func (c *Correction) CreateCorrection(ctx context.Context, userId int32, questionId string, newStatus domain.AnswerStatus) (string, error) {
	if !slices.Contains(domain.CorrectionAllowedNewStatuses, newStatus) {
		return "", errors.New("invalid new answer status")
	}
	correction := domain.Correction{
		ID:         domain.NewID(domain.ResourceTypeCorrection),
		QuestionId: questionId,
		UserId:     userId,
		NewStatus:  newStatus,
		UpdatedAt:  time.Now().UTC(),
	}
	err := c.questionStore.CreateCorrection(ctx, correction)
	if err != nil {
		return correction.ID, fmt.Errorf("failed to create correction: %w", err)
	}
	return correction.ID, nil
}

func (c *Correction) UpdateCorrectionNewStatus(ctx context.Context, correctionId string, newStatus domain.AnswerStatus) error {
	if !slices.Contains(domain.CorrectionAllowedNewStatuses, newStatus) {
		return errors.New("invalid new answer status")
	}
	return c.questionStore.UpdateCorrectionNewStatus(ctx, correctionId, newStatus)
}

func (c *Correction) UpdateCorrectionFeedback(ctx context.Context, correctionId string, feedback string) (domain.AnswerStatus, error) {
	return c.questionStore.UpdateCorrectionFeedback(ctx, correctionId, feedback)
}

func (c *Correction) FinalizeCorrection(ctx context.Context, correctionId string) error {
	return c.questionStore.FinalizeCorrection(ctx, correctionId)
}
