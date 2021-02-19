package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"
)

// Запускает воркер телеграма
func (a *MyApp) RunTelegramWorker() {
	var err error
	// Устанавливаем уровень журналирования событий функции
	log.SetLevel(a.logLevel)

	// Инициализируем подключение к телеге
	log.Info("Connecting to Telegram")
	a.bot, err = tgbotapi.NewBotAPI(a.config.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Authorized on account: %s (@%s)", a.bot.Self.FirstName, a.bot.Self.UserName)

	// Включаем глубокую отладку подключения к телеге, если в приложении включен максимальный дэбаг
	if a.logLevel == log.DebugLevel {
		a.bot.Debug = true
	}

	// Уведомляем коллег о запуске сервиса
	a.sendTelegramMessage(a.config.Telegram.ChatId, msgStatusConnected)

	// Дальше висим и слушаем все входящие сообщения
	a.listenTelegramMessages()

}

// Отправка сообщения в телеграм
func (a *MyApp) sendTelegramMessage(chatid int64, text string) {
	log.Infof("Telegram send [%d] %s", chatid, text)
	_, err := a.bot.Send(tgbotapi.NewMessage(chatid, text))
	if err != nil {
		log.Panic(err)
	}
}

// Получение и обработка входящих сообщений
func (a *MyApp) listenTelegramMessages() {
	log.Info("Starting telegram listener")
	// Устанавливаем время обновления канала сообщений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Инициализируем канал сообщений
	updates, err := a.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Optional: wait for updates and clear them if you don't want to handle
	// a large backlog of old messages
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	// Обрабатываем каждое полученное сообщение по каналу
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		chatId := "PRIVATE"
		if update.Message.Chat.Type != "private" {
			chatId = fmt.Sprintf("%s:%d", update.Message.Chat.Type, update.Message.Chat.ID)
		}

		log.Printf("Telegram read [%s][%s] %s", chatId, update.Message.From.UserName, update.Message.Text)

		// TODO вынести команды в отдельные функции
		// Обрабатываем команды
		if update.Message.IsCommand() {
			var reply string
			switch update.Message.Command() {
			case tgCommandStatus:
				reply = msgStatusAlive
			default:
				reply = msgUnknownCommand
			}
			a.sendTelegramMessage(update.Message.Chat.ID, reply)
		}
	}
}
