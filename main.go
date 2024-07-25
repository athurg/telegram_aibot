package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"
	"z/openai"
	"z/telegram"
)

var (
	flagBotToken    string
	flagOpenAiKey   string
	flagAiModel     string
	flagAdminUserId int64
)

func init() {
	flag.StringVar(&flagBotToken, "t", "", "Telegram bot token")
	flag.StringVar(&flagOpenAiKey, "k", "", "OpenAI key")
	flag.StringVar(&flagAiModel, "m", "gpt-4o-mini", "AI model")
	flag.Int64Var(&flagAdminUserId, "a", 0, "Telegram admin user id")
	flag.Parse()
}

func main() {
	if flagAdminUserId == 0 {
		log.Fatal("Admin user id is required.")
	}

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

			if update.Message.Chat.Type != "private" {
				_, err := tgBot.SendMessage(update.Message.Chat.Id, "Only private chat be supported.")
				if err != nil {
					log.Println("        Fail to reply message:", err)
				}
				continue
			}

			text := update.Message.TextValue()
			if update.Message.Chat.Id != flagAdminUserId {
				_, err := tgBot.SendMessage(update.Message.Chat.Id, fmt.Sprintf("Sorry, you (%d) are not allowed to chat with me.", update.Message.Chat.Id))
				if err != nil {
					log.Println("        Fail to reply message:", err)
				}
				continue
			}

			replyMessage, err := tgBot.SendMessage(update.Message.Chat.Id, "Connecting  AI ...")
			if err != nil {
				log.Println("        Fail to reply message:", err)
				continue
			}

			log.Printf("        Reply message: %+v", replyMessage)

			var replyText string
			lastReplyAt := time.Now()
			tokenCount, err := askAi(text, func(text string) {
				replyText += text

				//Buffered AI response to decrease the Telegram Bot API calls
				if time.Since(lastReplyAt) > time.Millisecond*500 {
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

			replyText += fmt.Sprintf("\n\nToken consume: %d", tokenCount)
			_, err = tgBot.EditMessageText(*replyMessage, replyText)
			if err != nil {
				log.Println("        Fail to reply message:", err)
			}
		}
	}
}

func askAi(question string, streamCallback func(string)) (int, error) {
	oaiCli := openai.New(flagOpenAiKey)
	oaiMessages := []openai.Message{
		openai.NewUserMessage(question),
	}
	var tokenCount int
	err := oaiCli.ChatCompletionsSteam(flagAiModel, oaiMessages, func(resp openai.ChatCompletionsResp) error {
		if resp.Usage != nil {
			tokenCount = resp.Usage.TotalTokens
			return nil
		}
		if len(resp.Choices) != 1 {
			return fmt.Errorf("Invalid AI response: %+v", resp)
		}

		if resp.Choices[0].Delta.Content == "" {
			return nil
		}

		streamCallback(resp.Choices[0].Delta.Content)
		return nil
	})

	return tokenCount, err
}
