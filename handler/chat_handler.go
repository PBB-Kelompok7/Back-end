package handler

import (
	"crowdfunding-minpro-alterra/modules/chat"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatUseCase chat.ChatUseCase
}

func NewChatHandler(chatUC chat.ChatUseCase) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUC,
	}
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	payload, err := chat.ParseChatPayload(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to parse request payload")
		return
	}

	response, err := h.chatUseCase.CompleteChat(payload)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to complete chat")
		return
	}

	// Ubah respons menjadi string
	responseString := string(response)

	// Pisahkan respons menjadi array objek chat
	chatObjects := strings.Split(responseString, "data:")

	// Buat slice untuk menyimpan objek chat dalam bentuk string
	var chatStrings []string

	// Uraikan setiap objek chat menjadi struktur data dan kirim hanya bagian "choices" sebagai respons JSON
	for i := 1; i < len(chatObjects)-1; i++ {
		obj := strings.TrimSpace(chatObjects[i])
		if obj == "" {
			continue
		}

		// Periksa apakah teks objek JSON yang valid sebelum mencoba mengurai
		if strings.HasPrefix(obj, "{") && strings.HasSuffix(obj, "}") {
			// Uraikan objek JSON
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(obj), &data); err != nil {
				c.String(http.StatusInternalServerError, "Failed to parse chat data")
				return
			}

			// Ambil bagian "choices" saja dari objek
			choices, ok := data["choices"].([]interface{})
			if !ok {
				c.String(http.StatusInternalServerError, "Failed to extract choices from chat data")
				return
			}

			// Kirim hanya bagian "choices" sebagai respons JSON tanpa kurung siku
			for _, choice := range choices {
				choiceJSON, err := json.Marshal(choice)
				if err != nil {
					c.String(http.StatusInternalServerError, "Failed to marshal chat choice")
					return
				}
				chatStrings = append(chatStrings, string(choiceJSON))
			}
		}
	}

	// Gabungkan semua objek chat menjadi satu string dengan pemisah koma
	responseJSON := "[" + strings.Join(chatStrings, ",") + "]"

	// Kirim respons JSON
	c.String(http.StatusOK, responseJSON)
}
