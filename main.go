package main

import (
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
}

func main() {
	var App MyApp
	var err error

	App.GetConfigYaml("conf.yml")

	log.SetLevel(App.logLevel)

	log.Info("Connecting to server...")

	// Connect to server
	App.imapClient, err = client.DialTLS(App.imapServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected")

	// Don't forget to logout
	defer App.imapClient.Logout()

	// Login
	if err := App.imapClient.Login(App.imapUsername, App.imapPassword); err != nil {
		log.Fatal(err)
	}
	log.Info("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- App.imapClient.List("", "*", mailboxes)
	}()

	log.Info("Mailboxes:")
	for m := range mailboxes {
		log.Info("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := App.imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Flags for INBOX:", mbox.Flags)

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	var uids []uint32
	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)

	go func() {
		uids, err = App.imapClient.Search(criteria)
		if err != nil {
			log.Error(err)
		}

		if len(uids) == 0 {
			log.Info("No new messages in go routine!")
			time.Sleep(10 * time.Second)
			return
		}

		log.Info("There are ", len(uids), " new messages")
		seqset := new(imap.SeqSet)
		seqset.AddNum(uids...)

		go func() {
			done <- App.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

	}()

	for msg := range messages {
		log.Info("* " + msg.Envelope.Subject)
		curSeq := new(imap.SeqSet)
		curSeq.AddNum(msg.SeqNum)
		err := App.imapClient.Store(curSeq, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	log.Info("Done!")
}
