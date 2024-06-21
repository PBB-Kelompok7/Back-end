package chat

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatPayload struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponse struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"`
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Choices           []ChatChoice `json:"choices"`
	SystemFingerprint string       `json:"system_fingerprint"`
}

// ChatChoice adalah struktur data untuk pilihan dalam respons chat.
type ChatChoice struct {
	Index        int          `json:"index"`
	Delta        DeltaContent `json:"delta"`
	Logprobs     interface{}  `json:"logprobs"`
	FinishReason interface{}  `json:"finish_reason"`
}

// DeltaContent adalah struktur data untuk konten delta dalam respons chat.
type DeltaContent struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}