package chat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ChatRepository interface {
	SendChatRequest(payload ChatPayload) ([]byte, error)
}

type chatRepository struct{}

func NewChatRepository() ChatRepository {
	return &chatRepository{}
}

func (r *chatRepository) SendChatRequest(payload ChatPayload) ([]byte, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post("https://wgpt-production.up.railway.app/v1/chat/completions", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
