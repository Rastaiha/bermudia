package service

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"log/slog"
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
		IsCorrect:  false,
		CreatedAt:  time.Now().UTC(),
	}
	create := false
	lowerFilename := strings.ToLower(answer.Filename.String)
	if strings.Contains(lowerFilename, "false") || strings.Contains(answer.TextContent.String, "false") || strings.Contains(answer.TextContent.String, "0") {
		correction.IsCorrect = false
		create = true
	} else if strings.Contains(lowerFilename, "true") || strings.Contains(answer.TextContent.String, "true") || strings.Contains(answer.TextContent.String, "1") {
		correction.IsCorrect = true
		create = true
	}
	if create {
		if err := c.questionStore.CreateCorrection(ctx, correction); err != nil {
			slog.Error("failed to create correction", slog.String("error", err.Error()))
		}
	}
	return create
}

func (c *Correction) CreateCorrection(userId int32, questionId string, isCorrect bool) (string, error) {
	correction := domain.Correction{
		ID:         domain.NewID(domain.ResourceTypeCorrection),
		QuestionId: questionId,
		UserId:     userId,
		IsCorrect:  isCorrect,
		CreatedAt:  time.Now().UTC(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := c.questionStore.CreateCorrection(ctx, correction)
	if err != nil {
		return correction.ID, fmt.Errorf("failed to create correction: %w", err)
	}
	return correction.ID, nil
}

func (c *Correction) UpdateCorrection(ctx context.Context, correctionId string, newIsCorrect bool) error {
	return c.questionStore.UpdateCorrection(ctx, correctionId, newIsCorrect)
}
