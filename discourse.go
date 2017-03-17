package area51bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	TopicCreatedEventType = "topic_created"
	PostCreatedEventType  = "post_created"
	PostEditedEventType   = "post_edited"
	EventHeader           = "X-Discourse-Event"
	InstanceHeader        = "X-Discourse-Instance"
)

type Notify func(chat int64)

// Post represents Discourse post entry from webhook payload
type Post struct {
	ID              int    `json:"id"`
	AuthorName      string `json:"name"`
	UserName        string `json:"username"`
	Number          int    `json:"post_number"`
	Type            int    `json:"post_type"`
	Preview         string `json:"cooked"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	TopicID         int    `json:"topic_id"`
	TopicSlug       string `json:"topic_slug"`
	DisplayUsername string `json:"display_username"`
	Admin           bool   `json:"admin"`
	Staff           bool   `json:"staff"`
	UserID          int    `json:"user_id"`
}

type Topic struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Visible   bool   `json:"visible"`
	UserID    int    `json:"user_id"`
	Slug      string `json:"slug"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
}

type NewTopicPayload struct {
	Topic    *Topic `json:"topic"`
	User     *User  `json:"user"`
	ForumURL string
}

func (p *NewTopicPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID)
	return fmt.Sprintf("%s (%s) создал новый топик <a href=\"%s\">\"%s\"</a> на форуме", p.User.Name, p.User.UserName, url, p.Topic.Title)
}

type NewPostPayload struct {
	Topic    *Topic `json:"topic"`
	Post     *Post  `json:"post"`
	User     *User  `json:"user"`
	ForumURL string
}

func (p *NewPostPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)
	return fmt.Sprintf("%s (%s) написал новый <a href=\"%s\">пост в \"%s\"</a>", p.User.Name, p.User.UserName, url, p.Topic.Title)
}

type EditedPostPayload struct {
	Topic    *Topic `json:"topic"`
	Post     *Post  `json:"post"`
	ForumURL string
}

func (p *EditedPostPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)
	return fmt.Sprintf("%s (%s) обновил <a href=\"%s\">пост в \"%s\"</a>", p.Post.AuthorName, p.Post.UserName, url, p.Topic.Title)
}

// HandleDiscourseEvent gets request detects vent type and depending of
// that type creates and sends Telegram message
func HandleDiscourseEvent(ctx context.Context, header http.Header, body []byte) (Notification, error) {
	e := header.Get(EventHeader)
	adr := header.Get(InstanceHeader)

	switch e {
	case TopicCreatedEventType:
		return handleCreatedTopicEvent(ctx, adr, body)
	case PostCreatedEventType:
		return handleCreatedPostEvent(ctx, adr, body)
	case PostEditedEventType:
		return handleEditedPostEvent(ctx, adr, body)
	default:
		return nil, nil
	}
}

func handleCreatedTopicEvent(ctx context.Context, adr string, body []byte) (Notification, error) {
	t := &NewTopicPayload{}
	err := json.Unmarshal(body, t)
	if err != nil {
		log.Errorf(ctx, "handleCreatedTopicEvent: %s", err.Error())
		return nil, err
	}
	t.ForumURL = adr

	NotifySubscribersByTheme(ctx, ThemeDiscourse, func(chat int64) {
		SendFormattedMessage(ctx, chat, t.Message(), HTMLFormatting)
	})

	return t, nil
}

func handleCreatedPostEvent(ctx context.Context, adr string, body []byte) (Notification, error) {
	t := &NewPostPayload{}
	err := json.Unmarshal(body, t)
	if err != nil {
		log.Errorf(ctx, "handleCreatedPostEvent: %s", err.Error())
		return nil, err
	}
	t.ForumURL = adr

	NotifySubscribersByTheme(ctx, ThemeDiscourse, func(chat int64) {
		SendFormattedMessage(ctx, chat, t.Message(), HTMLFormatting)
	})

	return t, nil
}

func handleEditedPostEvent(ctx context.Context, adr string, body []byte) (Notification, error) {
	t := &EditedPostPayload{}
	err := json.Unmarshal(body, t)
	if err != nil {
		log.Errorf(ctx, "handleCreatedTopicEvent: %s", err.Error())
		return nil, err
	}
	t.ForumURL = adr

	NotifySubscribersByTheme(ctx, ThemeDiscourse, func(chat int64) {
		SendFormattedMessage(ctx, chat, t.Message(), HTMLFormatting)
	})

	return t, nil
}
