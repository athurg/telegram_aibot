package main

import (
	"flag"
	"log"
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
			if update.Message != nil && update.Message.Text != nil {
				log.Printf("        Text: %s", update.Message.TextValue())
			}
		}
	}
}
