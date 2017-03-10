package area51bot

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
}

type NewTopicPayload struct {
	Topic Topic `json:"topic"`
	User  User  `json:"user"`
}

type NewPostPayload struct {
	Topic Topic `json:"topic"`
	Post  Post  `json:"post"`
	User  User  `json:"user"`
}

// HandleEvent gets request detects vent type and depending of
// that type creates and sends Telegram message
func HandleEvent(r *http.Request) error {
	e := r.Header.Get(EventHeader)
	m := &Message{Type: e}

	log.Printf("TADA %s", e)

	switch e {
	case TopicCreatedEventType:
		handleCreatedTopicEvent(m)
	case PostCreatedEventType:
		handleCreatedPostEvent(m)
	default:
		return nil
	}

	PostNotification(m)

	return nil
}

func handleCreatedTopicEvent(m *Message) error {
	f, err := os.Open("")
	if err != nil {
		return err
	}

	t := &NewTopicPayload{}
	err = json.NewDecoder(f).Decode(t)
	if err != nil {
		return err
	}

	m.Name = t.User.Name
	m.UserName = t.User.UserName
	m.Text = t.Topic.Title

	return nil
}

func handleCreatedPostEvent(m *Message) error {
	f, err := os.Open("")
	if err != nil {
		return err
	}

	t := &NewPostPayload{}
	err = json.NewDecoder(f).Decode(t)
	if err != nil {
		return err
	}

	m.Name = t.User.Name
	m.UserName = t.User.UserName
	m.Text = t.Post.Preview

	log.Println("TADA")

	return nil
}
