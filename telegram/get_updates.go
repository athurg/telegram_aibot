package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Chat struct {
	Id        int64
	Type      string
	Title     *string
	Username  *string
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	IsForum   *bool   `json:"is_forum"`
}

type Message struct {
	MessageId int64 `json:"message_id"`
	Date      int64
	Text      *string
	From      *User
	Chat      *Chat
}

func (m *Message) TextValue() string {
	if m.Text != nil {
		return *m.Text
	}
	return ""
}

type Update struct {
	UpdateId int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

func (b *Bot) GetUpdates(offset int64) ([]Update, error) {
	params := url.Values{
		"offset":  {strconv.FormatInt(offset, 10)},
		"timeout": {"10"},
	}

	resp, err := http.PostForm(ApiAddr+b.Token+"/getUpdates", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respInfo struct {
		Ok          bool
		Description string
		Result      []Update
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
