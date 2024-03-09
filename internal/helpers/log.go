package helpers

import (
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

// TeleContextAttrs returns a list of slog.Attr from tele.Context.
func TeleContextAttrs(c tele.Context) []any {
	if c == nil {
		return []any{}
	}

	return []any{
		teleUpdateAttr(c.Update()),
		teleUserAttr(c.Sender()),
		teleChatAttr(c.Chat()),
	}
}

func teleUpdateAttr(update tele.Update) slog.Attr {
	return slog.Group(
		"update",
		slog.Int("id", update.ID),
		slog.String("type", teleUpdateType(update)),
	)
}

func teleUserAttr(user *tele.User) slog.Attr {
	if user == nil {
		return slog.Group("user")
	}

	attrs := []any{
		slog.Int64("id", user.ID),
		slog.String("first_name", user.FirstName),
	}

	if user.LastName != "" {
		attrs = append(attrs, slog.String("last_name", user.LastName))
	}
	if user.Username != "" {
		attrs = append(attrs, slog.String("username", user.Username))
	}
	attrs = append(attrs, slog.Bool("is_bot", user.IsBot))

	return slog.Group("user", attrs...)
}

func teleChatAttr(chat *tele.Chat) slog.Attr {

	if chat == nil {
		return slog.Group("chat")
	}

	attrs := []any{
		slog.Int64("id", chat.ID),
		slog.String("type", string(chat.Type)),
	}

	if chat.Title != "" {
		attrs = append(attrs, slog.String("title", chat.Title))
	}
	if chat.Username != "" {
		attrs = append(attrs, slog.String("username", chat.Username))
	}

	return slog.Group("chat", attrs...)
}

func teleUpdateType(update tele.Update) string {
	switch {
	case update.Message != nil:
		return "message"
	case update.EditedMessage != nil:
		return "edited_message"
	case update.ChannelPost != nil:
		return "channel_post"
	case update.EditedChannelPost != nil:
		return "edited_channel_post"
	case update.Query != nil:
		return "inline_query"
	case update.InlineResult != nil:
		return "chosen_inline_result"
	case update.Callback != nil:
		return "callback_query"
	case update.ShippingQuery != nil:
		return "shipping_query"
	case update.PreCheckoutQuery != nil:
		return "pre_checkout_query"
	case update.Poll != nil:
		return "poll"
	case update.PollAnswer != nil:
		return "poll_answer"
	case update.MyChatMember != nil:
		return "my_chat_member"
	case update.ChatMember != nil:
		return "chat_member"
	case update.ChatJoinRequest != nil:
		return "chat_join_request"
	default:
		return "unknown"
	}
}
