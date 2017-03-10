package area51bot

import (
	"net/http"
	"fmt"
	"os"

	"github.com/kalambet/go-utils"
	"io/ioutil"
	"appengine_internal/memcache"
)

const (
	HookTypeParam  ="type"
	DiscourseValue ="discourse"
)

func SetupRouter() {
	http.HandleFunc("/health/", healthHandler)
	http.HandleFunc(fmt.Sprintf("/hook/%s/", os.Getenv("DISCOURSE_WEBHOOK")), hookHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All your base are belong to us!")
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	utils.PrintInColor(string(body), utils.Green)

	/*switch r.URL.Query().Get(HookTypeParam) {
	case DiscourseValue:
		discourse.HandleEvent(r.Body)
	}*/
}
