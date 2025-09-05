package service

import (
	"context"
	"fmt"
	"github.com/Rastaiha/bermudia/internal/config"
	"github.com/Rastaiha/bermudia/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Correction struct {
	cfg           config.Config
	bot           *bot.Bot
	questionStore domain.QuestionStore
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

func NewCorrection(cfg config.Config, bot *bot.Bot, questionStore domain.QuestionStore) *Correction {
	return &Correction{
		cfg:           cfg,
		bot:           bot,
		questionStore: questionStore,
	}
}

const (
	tagCB     = "tag|"
	correctCB = "correct|"
	wrongCB   = "wrong|"
	revertCB  = "revert|"
)

func prefix(cb string) bot.Middleware {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.CallbackQuery != nil {
				if data, ok := strings.CutPrefix(update.CallbackQuery.Data, cb); ok {
					update.CallbackQuery.Data = data
				}
			}
			next(ctx, bot, update)
		}
	}
}

func (c *Correction) Start() {
	c.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, tagCB, bot.MatchTypePrefix, c.handleTag, prefix(tagCB))
	c.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, correctCB, bot.MatchTypePrefix, c.handleCorrect, prefix(correctCB))
	c.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, wrongCB, bot.MatchTypePrefix, c.handleWrong, prefix(wrongCB))
	c.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, revertCB, bot.MatchTypePrefix, c.handleRevert, prefix(revertCB))
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.bot.Start(ctx)
	}()
}

func (c *Correction) Stop() {
	c.cancel()
	c.wg.Wait()
}

func (c *Correction) HandleNewAnswer(question domain.BookQuestion, answer domain.Answer) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if c.cfg.DevMode && c.autoCorrect(ctx, answer) {
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
		caption := fmt.Sprintf("#%d #%s\n\nمتن سؤال:\n%s", answer.UserID, answer.QuestionID, question.Text)
		if answer.FileID.Valid {
			_, err = c.bot.SendDocument(ctx, &bot.SendDocumentParams{
				ChatID:      c.cfg.DefaultCorrectionGroup,
				Document:    &models.InputFileString{Data: answer.FileID.String},
				Caption:     caption,
				ReplyMarkup: keyboard,
			})
		} else if utf8.RuneCount([]byte(answer.TextContent.String)) > 1024 {
			_, err = c.bot.SendDocument(ctx, &bot.SendDocumentParams{
				ChatID: c.cfg.DefaultCorrectionGroup,
				Document: &models.InputFileUpload{
					Filename: fmt.Sprintf("%d_%s.txt", answer.UserID, answer.QuestionID),
					Data:     strings.NewReader(answer.TextContent.String),
				},
				Caption:     caption,
				ReplyMarkup: keyboard,
			})
		} else {
			_, err = c.bot.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      c.cfg.DefaultCorrectionGroup,
				Text:        fmt.Sprintf("%s\n\nپاسخ کاربر\n%s", caption, answer.TextContent.String),
				ReplyMarkup: keyboard,
			})
		}
		if err != nil {
			slog.Error("error sending answer by bot", "error", err)
		}
	}()
}

func (c *Correction) autoCorrect(ctx context.Context, answer domain.Answer) bool {
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

func (c *Correction) createCorrection(userId int32, questionId string, isCorrect bool) (string, error) {
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

func (c *Correction) handleTag(ctx context.Context, b *bot.Bot, update *models.Update) {
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
	username := update.CallbackQuery.From.Username
	if username == "" {
		username = update.CallbackQuery.From.FirstName
	}
	suffix := fmt.Sprintf("\n\n✍️ @%s داره تصحیح میکنه", username)
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

func (c *Correction) handleCorrection(ctx context.Context, b *bot.Bot, update *models.Update, isCorrect bool) {
	var userId int32
	var questionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%d %s", &userId, &questionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}
	correctionId, err := c.createCorrection(userId, questionId, isCorrect)
	if err != nil {
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	c.showRevert(ctx, b, update, correctionId, isCorrect, false)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}

func (c *Correction) handleCorrect(ctx context.Context, b *bot.Bot, update *models.Update) {
	c.handleCorrection(ctx, b, update, true)
}

func (c *Correction) handleWrong(ctx context.Context, b *bot.Bot, update *models.Update) {
	c.handleCorrection(ctx, b, update, false)
}

func (c *Correction) showRevert(ctx context.Context, b *bot.Bot, update *models.Update, correctionId string, isCorrect bool, fromRevert bool) {
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

func (c *Correction) handleRevert(ctx context.Context, b *bot.Bot, update *models.Update) {
	var newIsCorrect bool
	var correctionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%t %s", &newIsCorrect, &correctionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = c.questionStore.UpdateCorrection(ctx, correctionId, newIsCorrect)
	if err != nil {
		err = fmt.Errorf("failed to update correction: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	c.showRevert(ctx, b, update, correctionId, newIsCorrect, true)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}
