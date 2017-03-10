package area51bot

import (
	"io"
	"encoding/json"
)

func HandleEvent(r io.ReadCloser) {
	var a = new(interface{})

	json.NewDecoder(r).Decode(a)

	PostNotification()
}
