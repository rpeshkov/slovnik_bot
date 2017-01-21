package main

import (
	"log"

	"os"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("SLOVNIK_BOT_ID"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// инициализируем канал, куда будут прилетать обновления от API
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := ""

		w, err := GetTranslations(update.Message.Text)

		if err != nil || len(w.word) <= 0 {
			messageText = "Specified word not found :("
		} else {
			messageText = w.String()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = "markdown"
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)

	}
}
