package chat

type ChatUseCase interface {
	CompleteChat(payload ChatPayload) ([]byte, error)
}

type chatUseCase struct {
	chatRepository ChatRepository
}

func NewChatUseCase(chatRepo ChatRepository) ChatUseCase {
	return &chatUseCase{
		chatRepository: chatRepo,
	}
}

func (uc *chatUseCase) CompleteChat(payload ChatPayload) ([]byte, error) {
	return uc.chatRepository.SendChatRequest(payload)
}
