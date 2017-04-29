package teledisq

const (
	TopicCreatedEventType = "topic_created"
	PostCreatedEventType  = "post_created"
	PostEditedEventType   = "post_edited"
	EventHeader           = "X-Discourse-Event"
	InstanceHeader        = "X-Discourse-Instance"
)

type Notify func(chat int64)

// Post represents Discourse post entry from web hook payload
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
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	Visible   bool         `json:"visible"`
	UserID    int          `json:"user_id"`
	Slug      string       `json:"slug"`
	Details   *PostDetails `json:"details"`
}

type PostDetails struct {
	CreatedBy *UserDetils `json:"created_by"`
}

type UserDetils struct {
	ID             int    `json:"id"`
	UserName       string `json:"username"`
	AvatarTemplate string `json:"avatar_template"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
}
