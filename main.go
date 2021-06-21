package main

import (
	"./config"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()

	bot, err := tgbotapi.NewBotAPI(conf.BotConfig.TelegramToken)
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

		resp, err := http.Post(conf.BotConfig.Url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			log.Fatalln(err)
		}

		var result map[string]map[string]string

		json.NewDecoder(resp.Body).Decode(&result)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.From.FirstName+" "+update.Message.From.LastName+"! "+result["compliment"]["compliment"])
		msg.ReplyToMessageID = update.Message.MessageID

		msgReport := tgbotapi.NewMessage(conf.BotConfig.ReportChatId, update.Message.From.FirstName+" "+update.Message.From.LastName+" used bot")

		bot.Send(msg)
		bot.Send(msgReport)
	}
}
