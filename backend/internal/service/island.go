package service

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
)

type Island struct {
	islandStore   domain.IslandStore
	questionStore domain.QuestionStore
	bot           *bot.Bot
}

func NewIsland(bot *bot.Bot, islandStore domain.IslandStore, questionStore domain.QuestionStore) *Island {
	return &Island{
		islandStore:   islandStore,
		questionStore: questionStore,
		bot:           bot,
	}
}

func (i *Island) GetIsland(ctx context.Context, userId int32, islandId string) (*domain.IslandContent, error) {
	rawContent, err := i.islandStore.GetByID(ctx, islandId)
	if err != nil {
		return nil, err
	}

	content := &domain.IslandContent{}
	for _, c := range rawContent.Components {
		if c.IFrame != nil {
			content.Components = append(content.Components, domain.IslandComponent{IFrame: c.IFrame})
			continue
		}
		if c.Question != nil {
			question, err := i.questionStore.GetQuestion(ctx, c.Question.QuestionID)
			if err != nil {
				return nil, err
			}
			userComponent, err := i.islandStore.GetOrCreateUserComponent(ctx, islandId, userId, c.ID, domain.ResourceTypeAnswer)
			if err != nil {
				return nil, err
			}
			answer, err := i.questionStore.GetOrCreateAnswer(ctx, userId, userComponent.ResourceID, question.ID)
			if err != nil {
				return nil, err
			}
			content.Components = append(content.Components, domain.IslandComponent{
				Input: &domain.IslandInput{
					ID:              answer.ID,
					Type:            question.InputType,
					Accept:          question.InputAccept,
					Description:     question.Text,
					SubmissionState: domain.GetSubmissionStateFromAnswer(answer),
				},
			})
			continue
		}
	}

	return content, nil
}

func (i *Island) SubmitAnswer(ctx context.Context, userId int32, answerId string, file io.ReadCloser, filename string) (*domain.SubmissionState, error) {
	// TODO: check player is in the island

	msg, err := i.bot.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: i.bot.ID(),
		Document: &models.InputFileUpload{
			Data:     file,
			Filename: filename,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file by bot sendDocument: %w", err)
	}
	if msg == nil || msg.Document == nil || msg.Document.FileID == "" {
		return nil, fmt.Errorf("failed to upload file by bot sendDocument: document is empty")
	}
	fileId := msg.Document.FileID

	answer, err := i.questionStore.SubmitAnswer(ctx, answerId, userId, fileId, filename)
	if err != nil {
		return nil, err
	}

	r := domain.GetSubmissionStateFromAnswer(answer)
	return &r, nil
}
