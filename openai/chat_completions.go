package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ChatCompletionsParam struct {
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
	Messages []Message `json:"messages"`
}

type ChatCompletionsResp struct {
	Id               string
	Object           string
	Created          int
	Model            string
	SystemFigerprint string `json:"system_fingerprint"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}
	Choices []struct {
		Index int
		Delta struct {
			Content string
		}
		Logprobs     any
		FinishReason string `json:"finish_reason"`
	}
}

type StreamCallback func(resp ChatCompletionsResp) error

func (c *Client) ChatCompletionsSteam(model string, messages []Message, callback StreamCallback) error {
	param := ChatCompletionsParam{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}

	paramsBytes, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ApiUrl+"/chat/completions", bytes.NewReader(paramsBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ChatCompletionsSteam failed: %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if line == "data: [DONE]" {
			break
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		var respInfo ChatCompletionsResp
		err := json.Unmarshal([]byte(line[6:]), &respInfo)
		if err != nil {
			return err
		}

		err = callback(respInfo)
		if err != nil {
			return err
		}
	}

	return nil
}
