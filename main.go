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

	for range time.NewTicker(2 * time.Second).C {
		App.searchNewMessages(criteria)
	}

	log.Info("Done!")
}
