package chat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ParseChatPayload(r *http.Request) (ChatPayload, error) {
	var payload ChatPayload
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return payload, err
	}
	err = json.Unmarshal(body, &payload)
	return payload, err
}
