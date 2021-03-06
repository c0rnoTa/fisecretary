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

// will be filled at buid phase
var gitHash, buildTime string

func main() {
	var App MyApp

	log.Info("fisecretary version ", gitHash, " build at ", buildTime)

	// Читаем конфиг
	App.GetConfigYaml(configFileName)

	// Устанавливаем уровень журналирования событий приложения
	log.SetLevel(App.logLevel)

	// Запускаем воркер работы с Telegram
	go App.RunTelegramWorker()

	// Запускаем воркер работы с IMAP
	if App.config.Imap.Enable {
		go func() {
			// Будет переподключаться, если разорвалось соединение
			for {
				App.RunImapWorker()
			}
		}()
	}

	// Запускаем подключение к Asterisk
	if App.config.Asterisk.Enable {
		go App.RunAsteriskWorker()
	}

	// TODO Сюда можно добавить проверку статуса подключений и их восстановление в цикле
	ch := make(chan bool)
	<-ch

	log.Info("Done!")
}
