package chat

import (
	"encoding/json"
)

// ResponseFormat adalah struktur untuk mengatur format respons yang diinginkan
type ResponseFormat struct {
	Delta        struct {
		Content string `json:"content"`
	} `json:"delta"`
	FinishReason interface{} `json:"finish_reason"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
}

// FormatResponse mengambil data respons dan mengembalikan data yang diformat sesuai kebutuhan
func FormatResponse(response []byte) ([]byte, error) {
	// Deklarasikan slice ResponseFormat untuk menyimpan objek-objek yang diformat
	var formattedResponse []ResponseFormat

	// Unmarshal respons JSON ke dalam slice ResponseFormat
	if err := json.Unmarshal(response, &formattedResponse); err != nil {
		return nil, err
	}

	// Return formatted response
	return json.Marshal(formattedResponse)
}
