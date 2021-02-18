package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

func (a *MyApp) sendTelegramMessage(chatid int64, text string) {
	_, err := a.bot.Send(tgbotapi.NewMessage(chatid, text))
	if err != nil {
		log.Panic(err)
	}
}

func (a *MyApp) receiveTelegramMessage() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatId := "PRIVATE"
		if update.Message.Chat.Type != "private" {
			chatId = fmt.Sprintf("%s:%d", update.Message.Chat.Type, update.Message.Chat.ID)
		}
		log.Printf("[%s][%s] %s", chatId, update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			var reply string
			switch update.Message.Command() {
			case "status":
				reply = "\u2705 Я в порядке!"
			default:
				reply = "Не понимаю, чего ты хочешь."
			}

			a.sendTelegramMessage(update.Message.Chat.ID, reply)
		}
	}
}
