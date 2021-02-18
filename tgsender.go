package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func (a *MyApp) sendMessage(text string) {
	_, err := a.bot.Send(tgbotapi.NewMessage(a.botChatId, text))
	if err != nil {
		log.Panic(err)
	}
}
