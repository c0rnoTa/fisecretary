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
		log.Fatal(err)
	}
	log.Info("IMAP Connected")

	// Don't forget to logout from IMAP server
	defer a.imapClient.Logout()

	// Login
	err = a.imapClient.Login(a.config.Imap.Username, a.config.Imap.Password)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Logged in as ", a.config.Imap.Username)

	// Выбираем папку INBOX на почтовом сервере
	log.Info("Select INBOX mailbox")
	_, err = a.imapClient.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}

		// Получаем UID-ы непрочитанных писем
		uids, err := a.imapClient.Search(criteria)
		if err != nil {
			log.Error(err)
		}
		// Если UID-ов нет, то новых писем нет
		if len(uids) == 0 {
			log.Info("No new messages yet.")
			return
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
				log.Fatal(err)
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
				log.Fatal(err)
			}
		}
	}

}
