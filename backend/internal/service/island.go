package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"io"
	"log/slog"
	"time"
)

type Island struct {
	bot                 *bot.Bot
	userStore           domain.UserStore
	islandStore         domain.IslandStore
	questionStore       domain.QuestionStore
	playerStore         domain.PlayerStore
	treasureStore       domain.TreasureStore
	gameStateStore      domain.GameStateStore
	onNewAnswer         NewAnswerCallback
	onNewPortableIsland NewPortableIslandCallback
	onHelpRequest       HelpRequestCallback
}

type NewAnswerCallback func(username string, territory string, question domain.BookQuestion, answer domain.Answer)

type NewPortableIslandCallback func(userId int32)

type HelpRequestCallback func(territory string, user *domain.User, question domain.BookQuestion) error

func NewIsland(bot *bot.Bot, userStore domain.UserStore, islandStore domain.IslandStore, questionStore domain.QuestionStore, playerStore domain.PlayerStore, treasureStore domain.TreasureStore, gameStateStore domain.GameStateStore) *Island {
	return &Island{
		bot:            bot,
		userStore:      userStore,
		islandStore:    islandStore,
		questionStore:  questionStore,
		playerStore:    playerStore,
		treasureStore:  treasureStore,
		gameStateStore: gameStateStore,
	}
}

func (i *Island) Start() {
	_ = func(now time.Time) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()

		if paused, err := i.gameStateStore.GetIsPaused(ctx); err != nil || paused {
			return
		}

		answers, err := i.questionStore.GetPendingAnswers(ctx, now.UTC().Add(-20*time.Minute))
		if err != nil {
			slog.Error("failed to get pending answers", slog.String("error", err.Error()))
			return
		}

		for _, a := range answers {
			user, err := i.userStore.Get(ctx, a.UserID)
			if err != nil {
				continue
			}
			question, err := i.questionStore.GetQuestion(ctx, a.QuestionID)
			if err != nil {
				continue
			}
			territory := ""
			if islandHeader, err := i.islandStore.GetIslandHeaderByBookIdAndUserId(ctx, question.BookID, user.ID); err == nil {
				if !islandHeader.FromPool {
					territory = islandHeader.TerritoryID
				}
			}
			i.onNewAnswer(user.Username, territory, question, a)
		}
	}
	//go func() {
	//	resend(time.Now())
	//	for now := range time.Tick(20 * time.Minute) {
	//		resend(now)
	//	}
	//}()
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
		content.Treasures = append(content.Treasures, domain.GetIslandTreasureOfUserTreasure(userTreasure, player.AtIsland == islandId))
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

func (i *Island) OnHelpRequest(f HelpRequestCallback) {
	i.onHelpRequest = f
}

func (i *Island) RequestHelpToAnswer(ctx context.Context, user *domain.User, questionId string) (string, error) {
	if user.MeetLink == "" {
		return "", domain.ErrMeetUnavailable
	}

	question, err := i.questionStore.GetQuestion(ctx, questionId)
	if err != nil {
		return "", err
	}
	answer, err := i.questionStore.GetAnswer(ctx, user.ID, questionId)
	if err != nil {
		return "", err
	}
	if err := domain.CheckRequestHelp(question, answer); err != nil {
		return "", err
	}

	if !answer.RequestedHelp {
		islandHeader, err := i.islandStore.GetIslandHeaderByBookIdAndUserId(ctx, question.BookID, user.ID)
		if err != nil {
			return "", err
		}
		err = i.onHelpRequest(islandHeader.TerritoryID, user, question)
		if err != nil {
			return "", fmt.Errorf("failed to call help request callback: %w", err)
		}
		err = i.questionStore.MarkHelpRequest(ctx, user.ID, questionId)
		if err != nil {
			return "", err
		}
	}

	return user.MeetLink, nil
}

func (i *Island) SetHelpState(ctx context.Context, userId int32, questionId string, state domain.HelpState) error {
	return i.questionStore.SetHelpState(ctx, userId, questionId, state)
}
