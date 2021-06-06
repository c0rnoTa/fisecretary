package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
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

	// Парсим аргументы команды запуска таймера
	duration, measure := TimerParseArgs(args)

	// Если вернулся 0, то в аргументах передан мусор, а значит запустить таймер нельзя
	if duration == 0 {
		log.Debug("Timer could not be started because invalid args were passed")
		a.sendTelegramMessage(update.Message.Chat.ID, msgTimerFailed)
		return
	}

	// Запустить таймер
	go TimerStart(a, update, duration, measure)

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
