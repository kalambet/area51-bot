package area51bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	TopicCreatedEventType = "topic_created"
	PostCreatedEventType  = "post_created"
	EventHeader           = "X-Discourse-Event"
	InstanceHeader        = "X-Discourse-Instance"
)

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
	Topic    Topic `json:"topic"`
	User     User  `json:"user"`
	ForumURL string
}

func (p *NewTopicPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID)
	return fmt.Sprintf("%s (%s) создал новый топик на форуме `%s`\n%s", p.User.Name, p.User.UserName, p.Topic.Title, url)
}

type NewPostPayload struct {
	Topic    Topic `json:"topic"`
	Post     Post  `json:"post"`
	User     User  `json:"user"`
	ForumURL string
}

func (p *NewPostPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)
	return fmt.Sprintf("%s (%s) написал новый пост на форуме %s\n\n%s", p.User.Name, p.User.UserName, url, p.Post.Preview)
}

// HandleEvent gets request detects vent type and depending of
// that type creates and sends Telegram message
func HandleEvent(r *http.Request) (Notification, error) {
	e := r.Header.Get(EventHeader)
	adr := r.Header.Get(InstanceHeader)

	switch e {
	case TopicCreatedEventType:
		return handleCreatedTopicEvent(adr, r.Body)
	case PostCreatedEventType:
		return handleCreatedPostEvent(adr, r.Body)
	default:
		return nil, nil
	}
}

func handleCreatedTopicEvent(adr string, r io.ReadCloser) (Notification, error) {
	t := &NewTopicPayload{ForumURL: adr}
	err := json.NewDecoder(r).Decode(r)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func handleCreatedPostEvent(adr string, r io.ReadCloser) (Notification, error) {
	t := &NewPostPayload{ForumURL: adr}
	err := json.NewDecoder(r).Decode(t)
	if err != nil {
		return nil, err
	}

	return t, nil
}
