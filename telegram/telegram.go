package telegram

const ApiAddr = "https://api.telegram.org/bot"

type Bot struct {
	Token string
}

func NewBot(token string) *Bot {
	return &Bot{Token: token}
}
