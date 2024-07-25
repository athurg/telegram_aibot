package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (b *Bot) SendMessage(chatId int64, text string) (*Message, error) {
	params := url.Values{
		"chat_id": {strconv.FormatInt(chatId, 10)},
		"text":    {text},
	}

	resp, err := http.PostForm(ApiAddr+b.Token+"/sendMessage", params)
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
