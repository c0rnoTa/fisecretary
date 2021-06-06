package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

type timerMeasureType struct {
	durationType time.Duration
	description  string
}

var TimerMeasure = map[string]timerMeasureType{
	timerSeconds: {time.Second, timerSecondsDescription},
	timerMinutes: {time.Minute, timerMinutesDescription},
	timerHours:   {time.Hour, timerHoursDescription},
}

var TimerDefaultMeasure = timerMinutes

// Простой таймер для напоминалки
func StartTimer(a *MyApp, update tgbotapi.Update, duration int, measure string) {
	log.SetLevel(a.logLevel)

	// По-умолчанию время ставится в минутах
	durationType := TimerMeasure[TimerDefaultMeasure].durationType
	measureDescription := TimerMeasure[TimerDefaultMeasure].description

	// Если в карте TimerMeasure есть установленное измерение времени
	// то меняем значение по-умолчанию для измерения времени таймера
	if _, ok := TimerMeasure[measure]; ok {
		// ок - это bool, выполнится, если в карте будет найден элемент
		durationType = TimerMeasure[measure].durationType
		measureDescription = TimerMeasure[measure].description
	}

	log.Debug("Timer started for ", duration, measureDescription)
	a.sendTelegramMessage(update.Message.Chat.ID, fmt.Sprintf(msgTimerStarted, duration, measureDescription))
	time.Sleep(time.Duration(duration) * durationType)

	// Если команда была с ответом на сообщение, то напоминает о сообщении, на которое был ответ
	replyto := update.Message.MessageID
	if update.Message.ReplyToMessage != nil {
		replyto = update.Message.ReplyToMessage.MessageID
	}
	a.replyTelegramMessage(update.Message.Chat.ID, replyto, fmt.Sprintf("@%s %s", update.Message.From.UserName, msgTimerFinished))
}
