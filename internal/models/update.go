package models

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

// Chat - chat model.
type Chat struct {
	Id       int64  `db:"chat_id"`
	Type     string `db:"chat_type"`
	Title    string `db:"chat_title"`
	Username string `db:"chat_username"`
}

// Member - chat member model.
type Member struct {
	Id        int64  `db:"member_id"`
	FirstName string `db:"member_first_name"`
	LastName  string `db:"member_last_name"`
	Username  string `db:"member_username"`
	IsBot     bool   `db:"member_is_bot"`
}

// Update - chat member update model.
type Update struct {
	Id        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Status    string    `db:"status"`
	Chat
	Member
}

// NewUpdate - creates new update model from tele.ChatMemberUpdate.
func NewUpdate(u *tele.ChatMemberUpdate) *Update {
	return &Update{
		Status: string(u.NewChatMember.Role),
		Chat: Chat{
			Id:       u.Chat.ID,
			Type:     string(u.Chat.Type),
			Title:    u.Chat.Title,
			Username: u.Chat.Username,
		},
		Member: Member{
			Id:        u.NewChatMember.User.ID,
			FirstName: u.NewChatMember.User.FirstName,
			LastName:  u.NewChatMember.User.LastName,
			Username:  u.NewChatMember.User.Username,
			IsBot:     u.NewChatMember.User.IsBot,
		},
	}
}
