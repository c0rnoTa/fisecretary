package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

// todo добавить комментов
// Простой таймер для напоминалки
func StartTimer(a *MyApp, update tgbotapi.Update, duration int, measure string) {
	log.SetLevel(a.logLevel)

	var durationType time.Duration
	switch measure {
	case timerSeconds:
		durationType = time.Second
	case timerHours:
		durationType = time.Hour
	case timerMinutes:
		durationType = time.Minute
	default:
		durationType = time.Minute
	}

	log.Debug("Timer started for ", duration, durationType)
	time.Sleep(time.Duration(duration) * durationType)
	// Если команда была с ответом на сообщение, то напоминает о сообщении, на которое был ответ
	replyto := update.Message.MessageID
	if update.Message.ReplyToMessage != nil {
		replyto = update.Message.ReplyToMessage.MessageID
	}
	a.replyTelegramMessage(update.Message.Chat.ID, replyto, msgTimerFinished)
}
