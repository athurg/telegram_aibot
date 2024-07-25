package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
	CanConnectToBusiness    bool   `json:"can_connect_to_business"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
}

func (b *Bot) GetMe() (*User, error) {
	resp, err := http.PostForm(ApiAddr+b.Token+"/getMe", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var respInfo struct {
		Ok          bool
		Description string
		Result      User
	}

	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return nil, err
	}

	if !respInfo.Ok {
		return nil, fmt.Errorf("telegram API error: %s", respInfo.Description)
	}

	return &respInfo.Result, nil
}
