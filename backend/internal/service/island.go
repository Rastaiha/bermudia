package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
)

type Island struct {
	bot           *bot.Bot
	islandStore   domain.IslandStore
	questionStore domain.QuestionStore
	playerStore   domain.PlayerStore
	treasureStore domain.TreasureStore
	onNewAnswer   func(question domain.BookQuestion, answer domain.Answer)
}

func NewIsland(bot *bot.Bot, islandStore domain.IslandStore, questionStore domain.QuestionStore, playerStore domain.PlayerStore, treasureStore domain.TreasureStore) *Island {
	return &Island{
		bot:           bot,
		islandStore:   islandStore,
		questionStore: questionStore,
		playerStore:   playerStore,
		treasureStore: treasureStore,
	}
}

func (i *Island) GetIsland(ctx context.Context, userId int32, islandId string) (*domain.IslandContent, error) {
	player, err := i.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	if err := domain.PlayerHasAccessToIsland(player, islandId); err != nil {
		return nil, err
	}
	territoryId, err := i.islandStore.GetTerritory(ctx, islandId)
	if err != nil {
		return nil, err
	}

	bookId, err := i.islandStore.GetBookOfIsland(ctx, islandId, userId)
	if errors.Is(err, domain.ErrNoBookAssignedFromPool) {
		bookId, err = i.islandStore.AssignBookToIslandFromPool(ctx, territoryId, islandId, userId)
	}
	if err != nil {
		return nil, err
	}
	book, err := i.islandStore.GetBook(ctx, bookId)
	if err != nil {
		return nil, err
	}

	content := &domain.IslandContent{}
	for _, c := range book.Components {
		if c.IFrame != nil {
			content.Components = append(content.Components, domain.IslandComponent{IFrame: c.IFrame})
			continue
		}
		if c.Question != nil {
			answer, err := i.questionStore.GetOrCreateAnswer(ctx, userId, c.Question.ID)
			if err != nil {
				return nil, err
			}
			content.Components = append(content.Components, domain.IslandComponent{
				Input: &domain.IslandInput{
					ID:              c.Question.ID,
					Type:            c.Question.InputType,
					Accept:          c.Question.InputAccept,
					Description:     c.Question.Text,
					SubmissionState: domain.GetSubmissionStateFromAnswer(answer),
				},
			})
			continue
		}
	}
	for _, t := range book.Treasures {
		userTreasure, err := i.treasureStore.GetOrCreateUserTreasure(ctx, userId, t.ID)
		if err != nil {
			return nil, err
		}
		content.Treasures = append(content.Treasures, domain.GetIslandTreasureOfUserTreasure(userTreasure))
	}

	return content, nil
}

func (i *Island) SubmitAnswer(ctx context.Context, userId int32, questionId string, file io.ReadCloser, filename string, textContent string) (*domain.SubmissionState, error) {
	player, err := i.playerStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	question, err := i.questionStore.GetQuestion(ctx, questionId)
	if err != nil {
		return nil, err
	}
	accessibleBook, err := i.islandStore.GetBookOfIsland(ctx, player.AtIsland, userId)
	if err != nil {
		if errors.Is(err, domain.ErrNoBookAssignedFromPool) {
			return nil, domain.ErrQuestionNotRelatedToIsland
		}
		return nil, err
	}
	if question.BookID != accessibleBook {
		return nil, domain.ErrQuestionNotRelatedToIsland
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

	answer, err := i.questionStore.SubmitAnswer(ctx, userId, questionId, fileId, filename, textContent)
	if err != nil {
		return nil, err
	}
	i.onNewAnswer(question, answer)

	r := domain.GetSubmissionStateFromAnswer(answer)
	return &r, nil
}

func (i *Island) OnNewAnswer(f func(question domain.BookQuestion, answer domain.Answer)) {
	i.onNewAnswer = f
}
