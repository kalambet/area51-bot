package area51bot

import (
	"fmt"
	"net/http"
	"os"
)

const (
	HookTypeParam  = "type"
	DiscourseValue = "discourse"
)

func SetupRouter() {
	http.HandleFunc("/health/", healthHandler)
	http.HandleFunc(fmt.Sprintf("/hook/%s/", os.Getenv("DISCOURSE_WEBHOOK")), hookHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All your base are belong to us!")
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get(HookTypeParam) {
	case DiscourseValue:
		HandleEvent(r)
	}
}
