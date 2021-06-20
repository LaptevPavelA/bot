package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	telegramToken, exists := os.LookupEnv("TELEGRAM_TOKEN")
	url, existsUrl := os.LookupEnv("URL")
	if !exists || !existsUrl {
		return
	}
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		fmt.Println(update.Message.From.LanguageCode)

		body, err := json.Marshal(map[string]interface{}{})

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]map[string]string

		json.NewDecoder(resp.Body).Decode(&result)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.From.FirstName+" "+update.Message.From.LastName+"! "+result["compliment"]["compliment"])
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
