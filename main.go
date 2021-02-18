package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type MyApp struct {
	logLevel     log.Level
	imapClient   *client.Client
	imapUsername string
	imapPassword string
	imapServer   string
	imapRefresh  time.Duration
	bot          *tgbotapi.BotAPI
	botToken     string
	botChatId    int64
}

func main() {
	var App MyApp
	var err error

	App.GetConfigYaml("conf.yml")

	log.SetLevel(App.logLevel)

	App.bot, err = tgbotapi.NewBotAPI(App.botToken)
	if err != nil {
		log.Fatal(err)
	}

	if App.logLevel == log.DebugLevel {
		App.bot.Debug = true
	}

	log.Infof("Authorized on account: %s (@%s)", App.bot.Self.FirstName, App.bot.Self.UserName)
	App.sendTelegramMessage(App.botChatId, "\u270c Привет! Я в сети :)")
	go App.receiveTelegramMessage()

	log.Info("Connecting to imap://", App.imapServer)

	// Connect to server
	App.imapClient, err = client.DialTLS(App.imapServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("IMAP Connected")

	// Don't forget to logout
	defer App.imapClient.Logout()

	// Login
	if err := App.imapClient.Login(App.imapUsername, App.imapPassword); err != nil {
		log.Fatal(err)
	}
	log.Info("Logged in as ", App.imapUsername)

	// Select INBOX
	log.Info("Select INBOX mailbox")
	_, err = App.imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	//log.Info("Flags for INBOX:", mbox.Flags)

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	for range time.NewTicker(App.imapRefresh * time.Second).C {
		err := App.imapClient.Noop()
		if err != nil {
			log.Fatal(err)
		}
		App.searchNewMails(criteria)
	}

	log.Info("Done!")
}
