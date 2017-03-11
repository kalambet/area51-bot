package area51bot

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI

func RunBot(ctx *appengine.Context) {
	host := appengine.DefaultVersionHostname(*ctx)
	url := fmt.Sprintf("https://%s/telegram/%s", host, os.Getenv("TELEGRAM_ENDPOINT"))

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_SECRET"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.SetWebhook(url)
}

func PostNotification() {
	log.Printf()
}
