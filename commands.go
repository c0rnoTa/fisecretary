package main

import (
	"fmt"
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
	case tgCommandWho:
		cmdWho(a, update)
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

// Поиск в справочнике
func cmdWho(a *MyApp, update tgbotapi.Update) {
	log.SetLevel(a.logLevel)

	if !a.config.Crm.Enable {
		log.Debug("CRM module disabled")
		a.sendTelegramMessage(update.Message.Chat.ID, msgWhoNoCRM)
		return
	}

	// Получаем аргументы команды
	args := strings.Fields(update.Message.CommandArguments())
	log.Debug("Telegram command is who with ", len(args), " arguments")

	if len(args) == 0 {
		log.Debug("Who command could not be started because no args were passed")
		a.sendTelegramMessage(update.Message.Chat.ID, msgWhoFailed)
		return
	}

	// Парсим аргументы обращения к who
	phoneNumber := args[0]

	msg := fmt.Sprintf(msgWhoReply, phoneNumber)
	callerName, err := a.getCrmName(phoneNumber)
	if err != nil {
		log.Error("Error in requesting CRM: ", err)
	} else {
		if callerName != "" {
			msg = fmt.Sprintf("%s\n"+msgCallPerson, msg, callerName)
		}
	}
	a.sendTelegramMessage(a.config.Telegram.ChatId, msg)
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
