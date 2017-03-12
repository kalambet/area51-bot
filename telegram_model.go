package teledisq

import (
	"os"
	"strings"
)

const (
	HTMLFormatting     = "HTML"
	ThemeDiscourse     = "discourse"
	CommandSendMessage = "sendMessage"
)

type ExternalNotification interface {
	Message() string
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type MessageEntities struct {
	Type string `json:"type"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type PhotoSize struct{}

type Message struct {
	ID             int64              `json:"message_id"`
	From           *TelegramUser      `json:"from"`
	Chat           *Chat              `json:"chat"`
	Text           string             `json:"text"`
	Entities       *[]MessageEntities `json:"entities"`
	NewChatMember  *TelegramUser      `json:"new_chat_member"`
	LeftChatMember *TelegramUser      `json:"left_chat_member"`
	NewChatTitle   string             `json:"new_chat_title"`
	NewChatPhoto   *[]PhotoSize       `json:"new_chat_photo"`
	Location       *Location          `json:"location"`
}

func (m *Message) IsCommand() bool {
	return m.Text != "" && m.Text[0] == '/'
}

func (m *Message) MentionedMe() bool {
	return m.Text != "" && strings.Contains(m.Text, os.Getenv("TELEGRAM_BOT_USERNAME"))
}

func (m *Message) IsNewChatMemberMessage() bool {
	return m.NewChatMember != nil
}

func (m *Message) IsNewChatTitleMessage() bool {
	return m.NewChatTitle != ""
}

func (m *Message) IsNewChatPhotoMessage() bool {
	return m.NewChatPhoto != nil
}

type Update struct {
	ID                int64    `json:"update_id"`
	Message           *Message `json:"message"`
	EditedMessage     *Message `json:"edited_message"`
	ChannelPost       *Message `json:"channel_post"`
	EditedChannelPost *Message `json:"edited_channel_post"`
}
