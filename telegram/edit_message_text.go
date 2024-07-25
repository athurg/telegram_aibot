package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (b *Bot) EditMessageText(msg Message, text string) (*Message, error) {
	if msg.Chat == nil {
		return nil, fmt.Errorf("chat is nil")
	}

	params := url.Values{
		"chat_id":    {strconv.FormatInt(msg.Chat.Id, 10)},
		"message_id": {strconv.FormatInt(msg.MessageId, 10)},
		"text":       {text},
	}

	resp, err := http.PostForm(ApiAddr+b.Token+"/editMessageText", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respInfo struct {
		Ok          bool
		Description string
		Result      *Message
	}

	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return nil, err
	}

	if !respInfo.Ok {
		return nil, fmt.Errorf("telegram API error: %s", respInfo.Description)
	}

	return respInfo.Result, nil
}
