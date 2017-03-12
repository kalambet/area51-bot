package teledisq

import (
	"fmt"
	"strings"
)

const (
	TopicCreatedEventType = "topic_created"
	PostCreatedEventType  = "post_created"
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
	return fmt.Sprintf("%s (%s) created new <a href=\"%s\">tpoic</a>:\n%s", p.User.Name, p.User.UserName, url, p.Topic.Title)
}

type NewPostPayload struct {
	Topic    *Topic `json:"topic"`
	Post     *Post  `json:"post"`
	User     *User  `json:"user"`
	ForumURL string
}

func (p *NewPostPayload) Message() string {
	url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)

	preview := p.Post.Preview
	if strings.Contains(preview, "<div") || strings.Contains(preview, "<blockquote") {
		preview = "but it's not possible for Telegram to do a preview"
	}

	return fmt.Sprintf("%s (%s) made new <a href=\"%s\">post</a>:\n\n%s", p.User.Name, p.User.UserName, url, preview)
}
