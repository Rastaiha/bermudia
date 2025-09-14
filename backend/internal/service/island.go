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
	bot                 *bot.Bot
	islandStore         domain.IslandStore
	questionStore       domain.QuestionStore
	playerStore         domain.PlayerStore
	treasureStore       domain.TreasureStore
	onNewAnswer         NewAnswerCallback
	onNewPortableIsland NewPortableIslandCallback
}

type NewAnswerCallback func(username string, territory string, question domain.BookQuestion, answer domain.Answer)

type NewPortableIslandCallback func(userId int32)

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
	isPortable, err := i.islandStore.IsIslandPortable(ctx, userId, islandId)
	if err != nil {
		return nil, err
	}
	if err := domain.CheckPlayerAccessToIslandContent(player, islandId, isPortable); err != nil {
		return nil, err
	}

	islandHeader, err := i.islandStore.GetIslandHeader(ctx, islandId)
	if err != nil {
		return nil, err
	}

	if domain.ShouldBeMadePortableOnAccess(islandHeader) {
		added, err := i.islandStore.AddPortableIsland(ctx, userId, islandId)
		if err != nil {
			return nil, err
		}
		if added {
			i.onNewPortableIsland(userId)
		}
	}

	bookId, err := i.islandStore.GetBookOfIsland(ctx, islandId, userId)
	if errors.Is(err, domain.ErrNoBookAssignedFromPool) {
		bookId, err = i.islandStore.AssignBookToIslandFromPool(ctx, islandHeader.TerritoryID, islandId, userId)
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
			question, err := i.questionStore.GetQuestion(ctx, c.Question.ID)
			if err != nil {
				return nil, err
			}
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
					SubmissionState: domain.GetSubmissionState(question, answer),
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

func (i *Island) SubmitAnswer(ctx context.Context, user *domain.User, questionId string, file io.ReadCloser, filename string, textContent string) (*domain.SubmissionState, error) {
	player, err := i.playerStore.Get(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	question, err := i.questionStore.GetQuestion(ctx, questionId)
	if err != nil {
		return nil, err
	}

	islandID := player.AtIsland
	territoryID := ""
	isPortable := false
	if islandHeader, err := i.islandStore.GetIslandHeaderByBookIdAndUserId(ctx, question.BookID, user.ID); err == nil {
		islandID = islandHeader.ID
		if !islandHeader.FromPool {
			territoryID = islandHeader.TerritoryID
		}
		isPortable, err = i.islandStore.IsIslandPortable(ctx, user.ID, islandHeader.ID)
		if err != nil {
			return nil, err
		}
	}
	if err := domain.CheckPlayerAccessToIslandContent(player, islandID, isPortable); err != nil {
		return nil, err
	}

	answer, err := i.questionStore.GetAnswer(ctx, user.ID, questionId)
	if err != nil {
		return nil, err
	}

	if err := domain.CheckSubmit(question, answer); err != nil {
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

	answer, err = i.questionStore.SubmitAnswer(ctx, user.ID, questionId, fileId, filename, textContent, answer.UpdatedAt)
	if err != nil {
		return nil, err
	}

	i.onNewAnswer(user.Username, territoryID, question, answer)

	r := domain.GetSubmissionState(question, answer)
	return &r, nil
}

func (i *Island) OnNewAnswer(f NewAnswerCallback) {
	i.onNewAnswer = f
}

func (i *Island) OnNewPortableIsland(f NewPortableIslandCallback) {
	i.onNewPortableIsland = f
}
