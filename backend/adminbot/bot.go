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
	"regexp"
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
	tagCB         = "tag|"
	correctCB     = "correct|"
	halfCorrectCB = "half|"
	wrongCB       = "wrong|"
	revertCB      = "revert|"
	finalizeCB    = "finalize|"
	iGoCB         = "iGo|"
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
	m.islandService.OnHelpRequest(m.HandleHelpRequest)

	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, tagCB, bot.MatchTypePrefix, m.handleTag, prefix(tagCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, correctCB, bot.MatchTypePrefix, m.handleCorrect, prefix(correctCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, halfCorrectCB, bot.MatchTypePrefix, m.handleHalfCorrect, prefix(halfCorrectCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, wrongCB, bot.MatchTypePrefix, m.handleWrong, prefix(wrongCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, revertCB, bot.MatchTypePrefix, m.handleRevert, prefix(revertCB))
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, finalizeCB, bot.MatchTypePrefix, m.handleFinalize, prefix(finalizeCB))
	m.bot.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, m.handleFeedback)
	m.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, iGoCB, bot.MatchTypePrefix, m.handleIGo, prefix(iGoCB))
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
		territory, group := m.getGroup(territory)
		caption := m.getMetaData(territory, username, question)
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

func (m *Bot) HandleHelpRequest(territory string, user *domain.User, question domain.BookQuestion) error {
	territory, group := m.getGroup(territory)
	msg := m.getMetaData(territory, user.Username, question)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := m.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: group,
		Text:   "☎️☎️☎️\n#درخواست_میت\n" + msg + fmt.Sprintf("\n\nMeet Link: %s", user.MeetLink),
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{
					Text:         "من جواب میدم",
					CallbackData: iGoCB,
				},
			}},
		},
	})
	return err
}

func (m *Bot) getGroup(territory string) (string, int64) {
	if territory == "" {
		territory = "challenge"
	}
	group := m.cfg.DefaultCorrectionGroup
	if g, ok := m.cfg.CorrectionGroups[territory]; ok {
		group = g
	}
	return territory, group
}

func (m *Bot) getMetaData(territory string, username string, question domain.BookQuestion) string {
	return fmt.Sprintf("#%s\nUser: #%s\nQuestion: #%s\n\nمتن سؤال:\n%s", territory, username, question.QuestionID, question.Text)
}

func (m *Bot) handleTag(ctx context.Context, b *bot.Bot, update *models.Update) {
	keyboard := models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{
				Text:         "غلط بود",
				CallbackData: wrongCB + update.CallbackQuery.Data,
			}},
			{{
				Text:         "نصفش درست بود",
				CallbackData: halfCorrectCB + update.CallbackQuery.Data,
			}},
			{{
				Text:         "کاملاً درست بود",
				CallbackData: correctCB + update.CallbackQuery.Data,
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

func (m *Bot) handleCorrection(ctx context.Context, b *bot.Bot, update *models.Update, newStatus domain.AnswerStatus) {
	var userId int32
	var questionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%d %s", &userId, &questionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}
	correctionId, err := m.correction.CreateCorrection(ctx, userId, questionId, newStatus)
	if err != nil {
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	m.showRevert(ctx, b, update, correctionId, newStatus, false)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}

func (m *Bot) handleCorrect(ctx context.Context, b *bot.Bot, update *models.Update) {
	m.handleCorrection(ctx, b, update, domain.AnswerStatusCorrect)
}

func (m *Bot) handleHalfCorrect(ctx context.Context, b *bot.Bot, update *models.Update) {
	m.handleCorrection(ctx, b, update, domain.AnswerStatusHalfCorrect)
}

func (m *Bot) handleWrong(ctx context.Context, b *bot.Bot, update *models.Update) {
	m.handleCorrection(ctx, b, update, domain.AnswerStatusWrong)
}

func statusToString(status domain.AnswerStatus) string {
	switch status {
	case domain.AnswerStatusWrong:
		return "غلط"
	case domain.AnswerStatusHalfCorrect:
		return "نیمه درست"
	case domain.AnswerStatusCorrect:
		return "کاملاً درست"
	default:
		return "unknown"
	}
}

func statusToEmoji(status domain.AnswerStatus) string {
	switch status {
	case domain.AnswerStatusWrong:
		return "🔴"
	case domain.AnswerStatusHalfCorrect:
		return "🟡"
	case domain.AnswerStatusCorrect:
		return "🟢"
	default:
		return "?"
	}
}

func revertKeyboard(correctionId string, currentNewStatus domain.AnswerStatus) models.InlineKeyboardMarkup {
	keyboard := models.InlineKeyboardMarkup{}
	for _, a := range domain.CorrectionAllowedNewStatuses {
		if currentNewStatus != a {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
				{
					Text:         fmt.Sprintf("تغییر نتیجه تصحیح به '%s'", statusToString(a)),
					CallbackData: revertCB + fmt.Sprintf("%d %s", a, correctionId),
				},
			})
		}
	}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []models.InlineKeyboardButton{
		{
			Text:         "ثبت نهایی نتیجه تصحیح",
			CallbackData: finalizeCB + fmt.Sprintf("%s", correctionId),
		},
	})
	return keyboard
}

func (m *Bot) showRevert(ctx context.Context, b *bot.Bot, update *models.Update, correctionId string, currentNewStatus domain.AnswerStatus, fromRevert bool) {
	keyboard := revertKeyboard(correctionId, currentNewStatus)

	suffix := fmt.Sprintf("```[ID]%s```\n", correctionId) + "✍️ با ریپلای به این پیام، برای دانش آموز یک متن بازخورد ثبت کنید.\n\n" +
		fmt.Sprintf("%s نتیجه تصحیح: *%s*", statusToEmoji(currentNewStatus), statusToString(currentNewStatus))
	if fromRevert {
		suffix = fmt.Sprintf("↩️%s نتیجه تصحیح به *%s* تغییر کرد", statusToEmoji(currentNewStatus), statusToString(currentNewStatus))
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
	var newStatus domain.AnswerStatus
	var correctionId string
	_, err := fmt.Sscanf(update.CallbackQuery.Data, "%d %s", &newStatus, &correctionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = m.correction.UpdateCorrectionNewStatus(ctx, correctionId, newStatus)
	if err != nil {
		err = fmt.Errorf("failed to update correction: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}

	m.showRevert(ctx, b, update, correctionId, newStatus, true)
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})
}

func (m *Bot) handleFinalize(ctx context.Context, b *bot.Bot, update *models.Update) {
	correctionId := update.CallbackQuery.Data
	err := m.correction.FinalizeCorrection(ctx, correctionId)
	if err != nil {
		err = fmt.Errorf("failed to parse callback query data: %w", err)
		_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
		slog.Error("failed to handle update", "error", err)
		return
	}
	suffix := "☑️ نتیجه نهایی تصحیح ثبت شد."
	suffix = "\n\n" + suffix
	if update.CallbackQuery.Message.Message.Document != nil {
		_, err = b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Caption:   update.CallbackQuery.Message.Message.Caption + suffix,
		})
	} else {
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text:      update.CallbackQuery.Message.Message.Text + suffix,
		})
	}
}

var correctionIdPattern = regexp.MustCompile("```\\[ID](\\w+)```")

func (m *Bot) handleFeedback(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" || update.Message.ReplyToMessage == nil {
		return
	}

	groups := correctionIdPattern.FindStringSubmatch(update.Message.ReplyToMessage.Text)
	if len(groups) < 2 {
		return
	}
	id := groups[1]
	if !domain.IdHasType(id, domain.ResourceTypeCorrection) {
		return
	}

	currentNewStatus, err := m.correction.UpdateCorrectionFeedback(ctx, id, update.Message.Text)
	if err != nil {
		slog.Error("failed to update correction feedback", "error", err)
		return
	}
	suffix := "🗒️ بازخورد شما ثبت شد."
	suffix = "\n\n" + suffix
	if update.Message.ReplyToMessage.Document != nil {
		_, err = b.EditMessageCaption(ctx, &bot.EditMessageCaptionParams{
			ChatID:      update.Message.ReplyToMessage.Chat.ID,
			MessageID:   update.Message.ReplyToMessage.ID,
			Caption:     update.Message.ReplyToMessage.Caption + suffix,
			ReplyMarkup: revertKeyboard(id, currentNewStatus),
		})
	} else {
		_, err = b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.Message.ReplyToMessage.Chat.ID,
			MessageID:   update.Message.ReplyToMessage.ID,
			Text:        update.Message.ReplyToMessage.Text + suffix,
			ReplyMarkup: revertKeyboard(id, currentNewStatus),
		})
	}
}

func (m *Bot) handleIGo(ctx context.Context, b *bot.Bot, update *models.Update) {
	username := update.CallbackQuery.From.FirstName
	if update.CallbackQuery.From.Username != "" {
		username = "@" + update.CallbackQuery.From.Username
	}
	suffix := fmt.Sprintf("🤙 %s جواب میده", username)
	suffix = "\n\n" + suffix
	_, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		Text:      update.CallbackQuery.Message.Message.Text + suffix,
	})
}
