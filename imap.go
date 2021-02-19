package main

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	log "github.com/sirupsen/logrus"
	"time"
)

// Запуск почтового воркера
func (a *MyApp) RunImapWorker() {
	var err error
	// Устанавливаем уровень журналирования событий функции
	log.SetLevel(a.logLevel)

	// Подключаемся к серверу IMAP
	log.Info("Connecting to imap://", a.config.Imap.Server)
	a.imapClient, err = client.DialTLS(a.config.Imap.Server, nil)
	if err != nil {
		log.Error("IMAP TLS connection returned error: ", err)
	}
	log.Info("IMAP Connected")

	// Don't forget to logout from IMAP server
	defer func() {
		err = a.imapClient.Logout()
		if err != nil {
			log.Error("IMAP Logout error: ", err)
		}
	}()

	// Login
	err = a.imapClient.Login(a.config.Imap.Username, a.config.Imap.Password)
	if err != nil {
		log.Error("IMAP login returned error: ", err)
	}
	log.Info("IMAP Logged in as ", a.config.Imap.Username)

	// Выбираем папку INBOX на почтовом сервере
	log.Info("Select INBOX mailbox")
	_, err = a.imapClient.Select("INBOX", false)
	if err != nil {
		log.Error("IMAP Mailbox folder select returned error: ", err)
	}

	// Дальше в бесконечном цикле ищем новые сообщения и увдомляем о них коллег
	a.ReadNewMail()
}

// Уведомляем о новых письмах
func (a *MyApp) ReadNewMail() {
	log.Info("Starting mailbox poller")

	// Установка критериев отбора писем в папке
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}

	// В бесконечном цыкле проверяем почтовый ящик на новые письма
	for range time.NewTicker(time.Duration(a.config.Imap.RefreshTimeout) * time.Second).C {
		// Чекаем новые письма
		err := a.imapClient.Noop()
		if err != nil {
			log.Error("IMAP Mailbox refresh returned error: ", err)
		}

		// Получаем UID-ы непрочитанных писем
		uids, err := a.imapClient.Search(criteria)
		if err != nil {
			log.Error("IMAP mail search returned error: ", err)
		}
		// Если UID-ов нет, то новых писем нет
		if len(uids) == 0 {
			log.Debug("No new messages yet.")
			continue
		}

		log.Info("There are ", len(uids), " new messages")
		seqset := new(imap.SeqSet)
		seqset.AddNum(uids...)

		// Инициализируем канал обработки полученных писем
		messages := make(chan *imap.Message, 10)
		// Отдельным потоком отгружаем найденные письма в канал
		go func() {
			err := a.imapClient.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
			if err != nil {
				log.Error("IMAP mail fetch error: ", err)
			}
		}()

		// Обрабатываем каждое новое письмо
		for msg := range messages {
			log.Info("* " + msg.Envelope.Subject)
			// Уведомляем коллег о полученном письме
			a.sendTelegramMessage(a.config.Telegram.ChatId, msg.Envelope.Subject)
			// Помечаем письмо как прочитанное
			curSeq := new(imap.SeqSet)
			curSeq.AddNum(msg.SeqNum)
			err := a.imapClient.Store(curSeq, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.SeenFlag}, nil)
			if err != nil {
				log.Error("IMAP mark mail as readed error: ", err)
			}
		}
	}

}
