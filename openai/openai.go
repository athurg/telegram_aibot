package openai

const ApiUrl = "https://api.openai.com/v1"

type Client struct {
	ApiKey string
}

func New(apiKey string) *Client {
	return &Client{
		ApiKey: apiKey,
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewSystemMessage(text string) Message {
	return Message{
		Role:    "system",
		Content: text,
	}
}

func NewUserMessage(text string) Message {
	return Message{
		Role:    "user",
		Content: text,
	}
}
func NewAssistantMessage(text string) Message {
	return Message{
		Role:    "assistant",
		Content: text,
	}
}
