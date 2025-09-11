package adminbot

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/Rastaiha/bermudia/internal/service"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Bot struct {
	cfg           config.Config
	bot           *bot.Bot
	islandService *service.Island
	correction    *service.Correction
	userStore     domain.UserStore
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewBot(cfg config.Config, b *bot.Bot, islandService *service.Island, correction *service.Correction, userStore domain.UserStore) *Bot {
	m := &Bot{
		cfg:           cfg,
		bot:           b,
		islandService: islandService,
		correction:    correction,
		userStore:     userStore,
	}

	return m
}

const (
	tagCB     = "tag|"
	correctCB = "correct|"
	wrongCB   = "wrong|"
	revertCB  = "revert|"
)

func prefix(cb string) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.CallbackQuery != nil {
				if data, ok := strings.CutPrefix(update.CallbackQuery.Data, cb); ok {
					update.CallbackQuery.Data = data
				}
			}
			next(ctx, b, update)
		}
	}
}

func (m *Bot) Start() {
	m.islandService.OnNewAnswer(m.HandleNewAnswer)

	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, tagCB, bot.MatchTypePrefix, m.handleTag, prefix(tagCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, correctCB, bot.MatchTypePrefix, m.handleCorrect, prefix(correctCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, wrongCB, bot.MatchTypePrefix, m.handleWrong, prefix(wrongCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, revertCB, bot.MatchTypePrefix, m.handleRevert, prefix(revertCB))
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		m.bot.Start(ctx)
	}()
}

func (m *Bot) Stop() {
	m.cancel()
	m.wg.Wait()
}

func (m *Bot) HandleNewAnswer(username string, territory string, question domain.BookQuestion, answer domain.Answer) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if m.cfg.DevMode && m.correction.AutoCorrect(ctx, answer) {
			return
		}
		keyboard := models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{
					Text:         "تصحیح می کنم",
					CallbackData: tagCB + fmt.Sprintf("%d %s", answer.UserID, answer.QuestionID),
				},
			}},
		}
		var err error
		if territory == "" {
			territory = "challenge"
		}
		group := m.cfg.DefaultCorrectionGroup
		if g, ok := m.cfg.CorrectionGroups[territory]; ok {
			group = g
		}
		caption := fmt.Sprintf("#%s\nUser: #%s\nQuestion: #%s\n\nمتن سؤال:\n%s", territory, username, answer.QuestionID, question.Text)
		if answer.FileID.Valid {
			_, err = m.bot.SendDocument(ctx, &bot.SendDocumentParams{
				ChatID:      group,
				Document:    &models.InputFileString{Data: answer.FileID.String},
				Caption:     caption,
				ReplyMarkup: keyboard,
			})
		} else if utf8.RuneCount([]byte(answer.TextContent.String)) > 1024 {
			_, err = m.bot.SendDocument(ctx, &bot.SendDocumentParams{
				ChatID: group,
				Document: &models.InputFileUpload{
					Filename: fmt.Sprintf("%d_%s.txt", answer.UserID, answer.QuestionID),
					Data:     strings.NewReader(answer.TextContent.String),
				},
				Caption:     caption,
				ReplyMarkup: keyboard,
			})
		} else {
			_, err = m.bot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      group,
				Text:        fmt.Sprintf("%s\n\nپاسخ کاربر:\n%s", caption, answer.TextContent.String),
				ReplyMarkup: keyboard,
			})
		}
		if err != nil {
			slog.Error("error sending answer by bot", "error", err)
		}
	}()
}

func (m *Bot) handleTag(ctx context.Context, b *bot.Bot, update *models.Update) {
	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{
				Text:         "درست بود",
				CallbackData: correctCB + update.CallbackQuery.Data,
			}},
			{{
				Text:         "غلط بود",
				CallbackData: wrongCB + update.CallbackQuery.Data,
			}},
		},
	}
	username := update.CallbackQuery.From.FirstName
	if update.CallbackQuery.From.Username != "" {
		username = "@" + update.CallbackQuery.From.Username
	}
	suffix := fmt.Sprintf("\n\n✏️ %s داره تصحیح میکنه", username)
	var err error
	if update.CallbackQuery.Message.Message.Document != nil {
		_, err = b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Caption:     update.CallbackQuery.Message.Message.Caption + suffix,
			ReplyMarkup: keyboard,
		})
	} else {
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        update.CallbackQuery.Message.Message.Text + suffix,
			ReplyMarkup: keyboard,
		})
	}
	if err != nil {
		return
	}
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
	return
}

func (m *Bot) handleCorrection(ctx context.Context, b *bot.Bot, update *models.Update, isCorrect bool) {
	var userId int32
	var questionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%d %s", &userId, &questionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}
	correctionId, err := m.correction.CreateCorrection(userId, questionId, isCorrect)
	if err != nil {
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	m.showRevert(ctx, b, update, correctionId, isCorrect, false)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}

func (m *Bot) handleCorrect(ctx context.Context, b *bot.Bot, update *models.Update) {
	m.handleCorrection(ctx, b, update, true)
}

func (m *Bot) handleWrong(ctx context.Context, b *bot.Bot, update *models.Update) {
	m.handleCorrection(ctx, b, update, false)
}

func (m *Bot) showRevert(ctx context.Context, b *bot.Bot, update *models.Update, correctionId string, isCorrect bool, fromRevert bool) {
	buttonText := "تغییر نتیجه تصحیح به 'درست'"
	if isCorrect {
		buttonText = "تغییر نتیجه تصحیح به 'غلط'"
	}
	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{
				Text:         buttonText,
				CallbackData: revertCB + fmt.Sprintf("%t %s", !isCorrect, correctionId),
			}},
		},
	}
	suffix := "❎ نتیجه تصحیح: *غلط*"
	if isCorrect {
		suffix = "✅ نتیجه تصحیح: *درست*"
	}
	if fromRevert {
		suffix = "↩️ نتیجه تصحیح به *غلط* تغییر کرد"
		if isCorrect {
			suffix = "↩️ نتیجه تصحیح به *درست* تغییر کرد"
		}
	}
	suffix = "\n\n" + suffix
	var err error
	if update.CallbackQuery.Message.Message.Document != nil {
		_, err = b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Caption:     update.CallbackQuery.Message.Message.Caption + suffix,
			ReplyMarkup: keyboard,
		})
	} else {
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        update.CallbackQuery.Message.Message.Text + suffix,
			ReplyMarkup: keyboard,
		})
	}
	if err != nil {
		slog.Error("failed to edit message on correction", "error", err)
		return
	}
}

func (m *Bot) handleRevert(ctx context.Context, b *bot.Bot, update *models.Update) {
	var newIsCorrect bool
	var correctionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%t %s", &newIsCorrect, &correctionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m.correction.UpdateCorrection(ctx2, correctionId, newIsCorrect)
	if err != nil {
		err = fmt.Errorf("failed to update correction: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	m.showRevert(ctx, b, update, correctionId, newIsCorrect, true)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}
