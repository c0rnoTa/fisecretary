package main

import (
	"github.com/emersion/go-imap/client"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ivahaev/amigo"
	log "github.com/sirupsen/logrus"
)

// Здесь все активные хэндлеры приложения
type MyApp struct {
	config     Config
	logLevel   log.Level
	imapClient *client.Client
	bot        *tgbotapi.BotAPI
	ami        *amigo.Amigo
}

func main() {
	var App MyApp

	// Читаем конфиг
	App.GetConfigYaml("conf.yml")

	// Устанавливаем уровень журналирования событий приложения
	log.SetLevel(App.logLevel)

	// Запускаем воркер работы с Telegram
	go App.RunTelegramWorker()

	// Запускаем воркер работы с IMAP
	go App.RunImapWorker()

	// Запускаем подключение к Asterisk
	go App.RunAsteriskWorker()

	// TODO Сюда можно добавить проверку статуса подключений и их восстановление в цикле
	ch := make(chan bool)
	<-ch

	log.Info("Done!")
}
