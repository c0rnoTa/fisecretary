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

	log.Info("Connecting to imap://", App.imapServer)

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

	uids, err := App.imapClient.Search(criteria)
	if err != nil {
		log.Error(err)
	}
	if len(uids) == 0 {
		log.Info("No new messages here.")
		time.Sleep(5 * time.Second)
		return
	}

	log.Info("There are ", len(uids), " new messages")
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	messages := make(chan *imap.Message, 10)

	go func() {
		err := App.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		if msg == nil {
			continue
		}
		log.Info("* " + msg.Envelope.Subject)
		curSeq := new(imap.SeqSet)
		curSeq.AddNum(msg.SeqNum)
		err := App.imapClient.Store(curSeq, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Done!")
}
