package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
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
func TimerStart(a *MyApp, update tgbotapi.Update, duration int, measure string) {
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

// Парсит параметры аргумента
func TimerParseArgs(args []string) (int, string) {

	// Формат первого аргумента
	// Может быть либо просто число, либо число с суффиксом
	r, _ := regexp.Compile(`(^\d+)(s|m|h)?`)

	// Первый аргумент - это число, счетчик
	duration := 0
	// Второй аргумент - это шаг счётчика, измерение. По-дефолту - минуты
	measure := TimerDefaultMeasure

	firstArg := r.FindStringSubmatch(args[0])
	//fmt.Printf("!!!!DEBUG!!! %d",len(firstArg))

	switch len(firstArg) {
	case 2:
		duration, _ = strconv.Atoi(firstArg[0])
	case 3:
		duration, _ = strconv.Atoi(firstArg[1])
		measure = firstArg[2]
	default:
		return duration, measure
	}

	if len(args) > 1 {
		measure = args[1]
	}
	return duration, measure
}
