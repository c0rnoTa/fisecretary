package main

import (
	"github.com/emersion/go-imap"
	log "github.com/sirupsen/logrus"
)

func (a *MyApp) searchNewMessages(criteria *imap.SearchCriteria) {
	uids, err := a.imapClient.Search(criteria)
	if err != nil {
		log.Error(err)
	}
	if len(uids) == 0 {
		log.Info("No new messages here.")
		return
	}

	log.Info("There are ", len(uids), " new messages")
	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	messages := make(chan *imap.Message, 10)

	go func() {
		err := a.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for msg := range messages {
		log.Info("* " + msg.Envelope.Subject)
		a.sendMessage(msg.Envelope.Subject)
		curSeq := new(imap.SeqSet)
		curSeq.AddNum(msg.SeqNum)
		err := a.imapClient.Store(curSeq, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.SeenFlag}, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}
