package area51bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	HookTypeParam  = "type"
	DiscourseValue = "discourse"
)

func SetupRouter() {
	http.HandleFunc("/", healthHandler)
	http.HandleFunc(fmt.Sprintf("/hook/%s/", os.Getenv("DISCOURSE_WEBHOOK")), hookHandler)
	http.HandleFunc(fmt.Sprintf("/telegram/%s/", os.Getenv("TELEGRAM_WEBHOOK")), telegramHandler)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "All your base are belong to us!")
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf(ctx, "hookHandler error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.URL.Query().Get(HookTypeParam) {
	case DiscourseValue:
		_, err = HandleDiscourseEvent(ctx, r.Header, body)
	default:
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	u := Update{}
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		log.Errorf(ctx, "Telegram update decoding error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	HandleTelegramUpdate(ctx, u)
	w.WriteHeader(http.StatusOK)
}
