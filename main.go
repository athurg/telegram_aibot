package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
	"z/openai"
	"z/telegram"
)

var flagBotToken string
var flagOpenAiKey string
var flagAiModel string

func init() {
	flag.StringVar(&flagBotToken, "t", "", "Telegram bot token")
	flag.StringVar(&flagOpenAiKey, "k", "", "OpenAI key")
	flag.StringVar(&flagAiModel, "m", "gpt-4o-mini", "AI model")
	flag.Parse()
}

func main() {
	log.Println("Starting bot...")
	runTelegramBot()
}

func runTelegramBot() {
	tgBot := telegram.NewBot(flagBotToken)
	botInfo, err := tgBot.GetMe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bot info: %+v", botInfo)

	updateOffset := int64(-1)
	for {
		log.Println("Getting updates...")
		updates, err := tgBot.GetUpdates(updateOffset)
		if err != nil {
			log.Fatal(err)
		}

		for i, update := range updates {
			updateOffset = update.UpdateId + 1

			updateBytes, _ := json.Marshal(update)
			log.Printf("[%4d/%4d]Update: %s", i+1, len(updates), string(updateBytes))

			if update.Message == nil || update.Message.Chat == nil || update.Message.Text == nil {
				log.Println("        Unsupport update")
				continue
			}

			text := update.Message.TextValue()
			if strings.HasPrefix(text, "/") {
				log.Println("        Unsupport command: ", text)
				continue
			}

			replyMessage, err := tgBot.SendMessage(update.Message.Chat.Id, "...")
			if err != nil {
				log.Println("        Fail to reply message:", err)
				continue
			}

			log.Printf("        Reply message: %+v", replyMessage)

			var replyText string
			lastReplyAt := time.Now()
			err = askAi(text, func(text string) {
				replyText += text

				//Buffered AI response to decrease the Telegram Bot API calls
				if time.Since(lastReplyAt) > time.Second {
					lastReplyAt = time.Now()
					_, err := tgBot.EditMessageText(*replyMessage, replyText)
					if err != nil {
						log.Println("        Fail to reply message:", err)
					}
				}
			})
			if err != nil {
				log.Println("        Error AI response:", err)
			}

			_, err = tgBot.EditMessageText(*replyMessage, replyText)
			if err != nil {
				log.Println("        Fail to reply message:", err)
			}
		}
	}
}

func askAi(question string, streamCallback func(string)) error {
	oaiCli := openai.New(flagOpenAiKey)
	oaiMessages := []openai.Message{
		openai.NewUserMessage(question),
	}
	err := oaiCli.ChatCompletionsSteam(flagAiModel, oaiMessages, func(resp openai.ChatCompletionsResp) error {
		if len(resp.Choices) != 1 {
			return fmt.Errorf("Invalid AI response: %+v", resp)
		}

		if resp.Choices[0].Delta.Content == "" {
			return nil
		}

		streamCallback(resp.Choices[0].Delta.Content)
		return nil
	})

	return err
}
