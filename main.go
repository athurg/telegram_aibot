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
}
