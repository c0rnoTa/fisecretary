package main

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.yandex.ru:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	log.Print(c.State())
	// Login
	if err := c.Login("username", "password"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	var uids []uint32
	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)

	go func() {
		uids, err = c.Search(criteria)
		if err != nil {
			log.Println(err)
		}

		if len(uids) == 0 {
			log.Print("No new messages in go routine!")
			time.Sleep(10 * time.Second)
			return
		}

		log.Print("There are ", len(uids), " new messages")
		seqset := new(imap.SeqSet)
		seqset.AddNum(uids...)

		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

	}()

	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
		curSeq := new(imap.SeqSet)
		curSeq.AddNum(msg.SeqNum)
		err := c.Store(curSeq, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	log.Println("Done!")
}
