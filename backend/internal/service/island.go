package service

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
	"log/slog"
	"strings"
	"time"
)

type Island struct {
	islandStore   domain.IslandStore
	questionStore domain.QuestionStore
	playerStore   domain.PlayerStore
	bot           *bot.Bot
}

func NewIsland(bot *bot.Bot, islandStore domain.IslandStore, questionStore domain.QuestionStore, playerStore domain.PlayerStore) *Island {
	return &Island{
		islandStore:   islandStore,
		questionStore: questionStore,
		playerStore:   playerStore,
		bot:           bot,
	}
}

func (i *Island) GetIsland(ctx context.Context, userId int32, islandId string) (*domain.IslandContent, error) {
	rawContent, territoryID, err := i.islandStore.GetByID(ctx, islandId)
	if err != nil {
		return nil, err
	}
	player, err := i.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	if err := domain.PlayerHasAccessToIsland(player, islandId); err != nil {
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
			answer, err := i.questionStore.GetOrCreateAnswer(ctx, userId, userComponent.ResourceID, question.ID, territoryID)
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

func (i *Island) SubmitAnswer(ctx context.Context, userId int32, answerId string, file io.ReadCloser, filename string, textContent string) (*domain.SubmissionState, error) {
	player, err := i.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	if err := i.islandStore.ResourceIsRelatedToIsland(ctx, userId, player.AtIsland, answerId); err != nil {
		return nil, err
	}
	if err := domain.PlayerHasAccessToIsland(player, player.AtIsland); err != nil {
		return nil, err
	}

	fileId := ""
	if file != nil {
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
		fileId = msg.Document.FileID
	}

	answer, err := i.questionStore.SubmitAnswer(ctx, answerId, userId, fileId, filename, textContent)
	if err != nil {
		return nil, err
	}

	// TODO: remove
	func() {
		correction := domain.Correction{
			ID:        domain.NewID(domain.ResourceTypeCorrection),
			AnswerID:  answer.ID,
			IsCorrect: false,
			CreatedAt: time.Now().UTC(),
		}
		create := false
		lowerFilename := strings.ToLower(filename)
		if strings.Contains(lowerFilename, "false") || strings.Contains(textContent, "false") || strings.Contains(textContent, "0") {
			correction.IsCorrect = false
			create = true
		} else if strings.Contains(lowerFilename, "true") || strings.Contains(textContent, "true") || strings.Contains(textContent, "1") {
			correction.IsCorrect = true
			create = true
		}
		if create {
			if err := i.questionStore.CreateCorrection(ctx, correction); err != nil {
				slog.Error("failed to create correction", err)
			}
		}
	}()

	r := domain.GetSubmissionStateFromAnswer(answer)
	return &r, nil
}
