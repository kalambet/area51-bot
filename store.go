package teledisq

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const (
	SubscriptionKind = "Subscription"
)

type Subscription struct {
	ChatID int64
	Theme  string
}

func SubscribeChat(ctx context.Context, chat int64, theme string) (bool, error) {
	key := datastore.NewKey(ctx, SubscriptionKind, "", chat, nil)

	s := Subscription{}
	err := datastore.Get(ctx, key, &s)
	if err == datastore.ErrNoSuchEntity {
		_, err := datastore.Put(ctx, key, &Subscription{ChatID: chat, Theme: theme})
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func UnsubscribeChat(ctx context.Context, chat int64, _ string) (bool, error) {
	key := datastore.NewKey(ctx, SubscriptionKind, "", chat, nil)

	err := datastore.Delete(ctx, key)
	if err == datastore.ErrNoSuchEntity {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func NotifySubscribersByTheme(ctx context.Context, theme string, notify Notify) {
	q := datastore.NewQuery(SubscriptionKind).Filter("Theme=", theme)
	for t := q.Run(ctx); ; {
		var s Subscription
		_, err := t.Next(&s)
		if err == datastore.Done {
			break
		}

		if err != nil {
			log.Errorf(ctx, "Problem on subscribers query %s", err.Error())
			return
		}

		// Send notification
		notify(s.ChatID)
	}
}
