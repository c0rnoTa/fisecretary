package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Роутер обработчкиа команд
func RouteCommands(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)
	log.Debug("Processing telegram command '", update.Message.Command(), "'")

	// Вызываем функцию обработки команды в зависимости от переданной команды
	switch update.Message.Command() {
	case tgCommandStatus:
		cmdStatus(a, update)
	default:
		cmdUnknown(a, update)
	}
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
