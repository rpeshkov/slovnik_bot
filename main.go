package main

import (
	"log"

	"os"

	"github.com/rpeshkov/slovnik_bot/slovnik"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("SLOVNIK_BOT_ID"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		messageText := ""

		lang := slovnik.DetectLanguage(update.Message.Text)
		log.Println(lang)

		w, err := slovnik.GetTranslations(update.Message.Text, lang)

		if err != nil || len(w.Word) <= 0 {
			messageText = "Specified word not found :("
		} else {
			messageText = w.String()
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)

		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)

	}
}
