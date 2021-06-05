package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// Роутер обработчкиа команд
func RouteCommands(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)
	log.Debug("Processing telegram command '", update.Message.Command(), "'")

	// Вызываем функцию обработки команды в зависимости от переданной команды
	switch update.Message.Command() {
	case tgCommandStatus:
		cmdStatus(a, update)
	case tgCommandTimer:
		cmdTimer(a, update)
	default:
		cmdUnknown(a, update)
	}
}

// TODO: переделать на карту
// Напоминалка
func cmdTimer(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)

	// Получаем аргументы команды
	args := strings.Fields(update.Message.CommandArguments())
	log.Debug("Telegram command is timer with ", len(args), " arguments")

	if len(args) == 0 {
		log.Debug("Timer could not be started because no args were passed")
		a.sendTelegramMessage(update.Message.Chat.ID, msgTimerFailed)
		return
	}

	// Первый аргумент - это число, счетчик
	duration, err := strconv.Atoi(args[0])
	if err != nil {
		log.Debug("Timer could not be started because ", err)
		a.sendTelegramMessage(update.Message.Chat.ID, msgTimerFailed)
		return
	}
	// Второй аргумент - это шаг счётчика, измерение. По-дефолту - минуты
	measure := timerMinutes
	measureDescription := timerMinutesDescription
	if len(args) > 1 {
		measure = args[1]
	}
	switch measure {
	case timerSeconds:
		measureDescription = timerSecondsDescription
	case timerMinutes:
		measureDescription = timerMinutesDescription
	case timerHours:
		measureDescription = timerHoursDescription
	}
	// Запустить таймер
	go StartTimer(a, update, duration, measure)
	a.sendTelegramMessage(update.Message.Chat.ID, fmt.Sprintf(msgTimerStarted, duration, measureDescription))
}

// Обработчик команд по-умолчанию. Если команда не известна
func cmdUnknown(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)
	log.Debug("Telegram command is unknown")
	// Просто отвечаем, что не понимаем, что от нас хотят
	a.sendTelegramMessage(update.Message.Chat.ID, msgUnknownCommand)
}

// Проверка статуса работы сервиса бота
func cmdStatus(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)
	log.Debug("Telegram command is status check")
	a.sendTelegramMessage(update.Message.Chat.ID, msgStatusAlive)
}
