package teledisq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// DiscoursePayload describe the payload of Discourse eventx.
// Depending of event type it cotains different data set
// See your Discourse events requests for more information
type DiscoursePayload struct {
	Topic    *Topic `json:"topic"`
	Post     *Post  `json:"post"`
	User     *User  `json:"user"`
	ForumURL string
	Event    string
}

// Message generates appropriate notification message based on event type
func (p *DiscoursePayload) Message() string {
	switch p.Event {
	case PostCreatedEventType:
		url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)
		return fmt.Sprintf("%s (%s) написал новый <a href=\"%s\">пост в \"%s\"</a>", p.User.Name, p.User.UserName, url, p.Topic.Title)
	case PostEditedEventType:
		url := fmt.Sprintf("%s/t/%s/%d/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID, p.Post.ID)
		return fmt.Sprintf("%s (%s) обновил <a href=\"%s\">пост в \"%s\"</a>", p.Post.AuthorName, p.Post.UserName, url, p.Topic.Title)
	case TopicCreatedEventType:
		url := fmt.Sprintf("%s/t/%s/%d", p.ForumURL, p.Topic.Slug, p.Topic.ID)
		return fmt.Sprintf("%s (%s) создал новый топик <a href=\"%s\">\"%s\"</a> на форуме", p.User.Name, p.User.UserName, url, p.Topic.Title)
	}

	return ""

}

// HandleDiscourseEvent gets request detects event type and depending of
// that type creates and sends Telegram message
func HandleDiscourseEvent(ctx context.Context, header http.Header, body []byte) (string, error) {
	e := header.Get(EventHeader)
	forumURL := header.Get(InstanceHeader)

	p := DiscoursePayload{ForumURL: forumURL, Event: e}
	err := json.Unmarshal(body, &p)
	if err != nil {
		log.Errorf(ctx, "HandleDiscourseEvent: %s", err.Error())
		return "", err
	}

	return p.Message(), nil
}
