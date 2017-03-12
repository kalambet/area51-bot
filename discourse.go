package teledisq

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

// HandleDiscourseEvent gets request detects vent type and depending of
// that type creates and sends Telegram message
func HandleDiscourseEvent(ctx context.Context, header http.Header, body []byte) (ExternalNotification, error) {
	e := header.Get(EventHeader)
	forumURL := header.Get(InstanceHeader)

	switch e {
	case TopicCreatedEventType:
		return handleCreatedTopicEvent(ctx, forumURL, body)
	case PostCreatedEventType:
		return handleCreatedPostEvent(ctx, forumURL, body)
	default:
		return nil, nil
	}
}

func handleCreatedTopicEvent(ctx context.Context, forumURL string, body []byte) (ExternalNotification, error) {
	t := &NewTopicPayload{ForumURL: forumURL}
	err := json.Unmarshal(body, t)
	if err != nil {
		log.Errorf(ctx, "handleCreatedTopicEvent: %s", err.Error())
		return nil, err
	}

	return t, nil
}

func handleCreatedPostEvent(ctx context.Context, adr string, body []byte) (ExternalNotification, error) {
	t := &NewPostPayload{ForumURL: adr}
	err := json.Unmarshal(body, t)
	if err != nil {
		log.Errorf(ctx, "handleCreatedPostEvent: %s", err.Error())
		return nil, err
	}

	return t, nil
}
