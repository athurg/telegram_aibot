package main

import (
	"flag"
	"log"
	"time"
	"z/telegram"
)

var flagBotToken string

func init() {
	flag.StringVar(&flagBotToken, "t", "", "Telegram bot token")
	flag.Parse()
}

func main() {
	log.Println("Starting bot...")

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

			log.Printf("[%4d/%4d]Update: %+v", i+1, len(updates), update)
			if update.Message == nil {
				log.Println("        No message")
				continue
			}

			if update.Message.Chat == nil {
				log.Println("        No chat")
				continue
			}

			if update.Message.Text == nil {
				log.Println("        No text")
				continue
			}

			log.Printf("        [TEXT]: %s", update.Message.TextValue())

			replyMessage, err := tgBot.SendMessage(update.Message.Chat.Id, update.Message.TextValue())
			if err != nil {
				log.Println("        Error reply message:", err)
				continue
			}

			log.Printf("        Reply message: %+v", replyMessage)

			time.Sleep(1 * time.Second)
			editedMessage, err := tgBot.EditMessageText(*replyMessage, update.Message.TextValue()+"(Edited)")
			if err != nil {
				log.Println("        Error edit message:", err)
				continue
			}
			log.Printf("        Edited message: %+v", editedMessage)
		}
	}
}
