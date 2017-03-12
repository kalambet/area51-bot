package area51bot

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func HandleTelegramUpdate(ctx context.Context, u Update) {
	m := &Message{}
	if u.EditedMessage != nil {
		m = u.EditedMessage
	} else {
		m = u.Message
	}

	log.Infof(ctx, "Update message content: %+v", u)

	if m.IsCommand() {
		handleCommand(ctx, m)
	} else if strings.Contains(m.LeftChatMember.UserName, os.Getenv("TELEGRAM_BOT_USERNAME")) {
		handleRemoving(ctx, m)
	}
}

func handleRemoving(ctx context.Context, m *Message) {
	success, err := UnsubscribeChat(ctx, m.Chat.ID, ThemeDiscourse)
	if err != nil {
		log.Errorf(ctx, "Problem removing subscriotion for chat %d: %s", m.Chat.ID, err)
		return
	}

	if !success {
		log.Errorf(ctx, "Chat %d was not subscribed", m.Chat.ID)
		return
	}
}

func handleCommand(ctx context.Context, m *Message) {
	if m == nil {
		return
	}

	log.Infof(ctx, "Command: %s", m.Text)

	commands := strings.Split(m.Text, " ")
	if len(commands) == 0 || commands[0] != "/area51" {
		return
	}

	if strings.Contains(strings.ToUpper(m.Text), strings.ToUpper("не рассказывай нам про форум")) {
		SendMessage(ctx, m.Chat.ID, "чёт не ясно, что надо 😗")
	} else if strings.Contains(strings.ToUpper(m.Text), strings.ToUpper("рассказывай нам про форум")) {
		success, err := SubscribeChat(ctx, m.Chat.ID, ThemeDiscourse)
		if err != nil {
			log.Errorf(ctx, "Problem interacting with Datastore: %s", err.Error())
			SendMessage(ctx, m.Chat.ID, "неа")
			return
		}

		if success {
			SendMessage(ctx, m.Chat.ID, "нну ок")
		} else {
			SendMessage(ctx, m.Chat.ID, "так я же уже рассказываю вам про форум")
		}
	} else {
		SendMessage(ctx, m.Chat.ID, "¯\\_(ツ)_/¯")
		return
	}
}

func SendMessage(ctx context.Context, chat int64, text string) {
	SendFormattedMessage(ctx, chat, text, "")
}

func SendFormattedMessage(ctx context.Context, chat int64, text string, format string) {
	payload := make(url.Values)

	payload.Add("chat_id", fmt.Sprintf("%d", chat))
	payload.Add("text", sanitizeHTMLInput(text))

	if format != "" {
		payload.Add("parse_mode", format)
	}

	makeRequest(ctx, CommandSendMessage, payload)
}

func makeRequest(ctx context.Context, cmd string, data url.Values) {
	c := urlfetch.Client(ctx)
	if c == nil {
		log.Errorf(ctx, "Can't create AppEngine urlfetch Client")
		return
	}

	address := fmt.Sprintf("https://api.telegram.org/bot%s/%s", os.Getenv("TELEGRAM_SECRET"), cmd)

	// Always add the telegram method we use to POST
	data.Add("method", cmd)

	resp, err := c.PostForm(address, data)
	if err != nil {
		log.Errorf(ctx, "Fail to make send message request %s. Payload: %#v", cmd, data)
		return
	}

	if resp.StatusCode > 201 {
		log.Errorf(ctx, "Bad send message request for '%s'.\nStatus: %s\nPayload: %#v", cmd, resp.Status, data)
		return
	}
}

func sanitizeHTMLInput(text string) string {

	text = strings.Replace(text, "<p>", " ", -1)
	text = strings.Replace(text, "</p>", " ", -1)
	text = strings.Replace(text, "\\\"", "\"", -1)

	return text
}
