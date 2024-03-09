package notifier

import (
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v3"

	"joinfive-bot/internal/models"
)

// TelegramSingleChat is a notifier implementation for a single telegram recipient.
type TelegramSingleChat struct {
	to  tele.Recipient
	bot *tele.Bot
}

// NewTelegramSingleChat - TelegramSingleChat constructor.
func NewTelegramSingleChat(chatID int64, bot *tele.Bot) *TelegramSingleChat {
	slog.Info("notifier: created",
		slog.String("type", "TelegramSingleChat"),
		slog.Int64("recipient", chatID),
	)
	return &TelegramSingleChat{
		to:  &tele.Chat{ID: chatID},
		bot: bot,
	}
}

// Notify sends a message to the recipient.
func (t *TelegramSingleChat) Notify(update *models.Update) {
	_, err := t.bot.Send(t.to, t.fmtMsg(update), tele.ModeHTML, tele.NoPreview)
	if err != nil {
		slog.Error(fmt.Sprintf("notify: failed to send message: %v", err))
		return
	}

	slog.Info(
		"notify: message sent",
		slog.String("recipient", t.to.Recipient()),
	)
}

// fmtMsg formats the message to be sent.
func (t *TelegramSingleChat) fmtMsg(update *models.Update) string {
	return fmt.Sprintf(
		"%s %s → %s",
		t.fmtStatus(update.Status),
		t.fmtMember(update.Member),
		t.fmtChat(update.Chat),
	)
}

func (t *TelegramSingleChat) fmtMember(member models.Member) string {
	result := member.FirstName
	if member.LastName != "" {
		result += " " + member.LastName
	}
	result = fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, member.Id, result)
	if member.Username != "" {
		result = fmt.Sprintf("%s (@%s)", result, member.Username)
	}
	return result
}

func (t *TelegramSingleChat) fmtChat(chat models.Chat) string {
	result := chat.Title
	if chat.Username != "" {
		result = fmt.Sprintf(`<a href="https://t.me/%s">%s</a>`, chat.Username, result)
	}
	return result
}

func (t *TelegramSingleChat) fmtStatus(status string) string {
	result := ""
	switch status {
	case "creator", "administrator":
		result = "🎩"
	case "member":
		result = "🟢"
	case "left":
		result = "🔴"
	case "restricted":
		result = "❗"
	case "kicked":
		result = "❌"
	default:
		result = status + ":"
	}

	return result
}
