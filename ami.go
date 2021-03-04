package main

import (
	"fmt"
	"github.com/ivahaev/amigo"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (a *MyApp) RunAsteriskWorker() {
	// Устанавливаем уровень журналирования событий функции
	log.SetLevel(a.logLevel)
	log.Infof("Connecting to AMI %s:%d", a.config.Asterisk.Host, a.config.Asterisk.Port)
	settings := &amigo.Settings{
		Username: a.config.Asterisk.Username,
		Port:     strconv.Itoa(a.config.Asterisk.Port),
		Password: a.config.Asterisk.Password,
		Host:     a.config.Asterisk.Host,
	}
	a.ami = amigo.New(settings)

	a.ami.Connect()

	a.ami.On("connect", func(message string) {
		log.Info("Connected to PBX: ", message)
	})

	a.ami.On("error", func(message string) {
		amiConn := fmt.Sprintf("%s:%s@%s:%d", a.config.Asterisk.Username, a.config.Asterisk.Password, a.config.Asterisk.Host, a.config.Asterisk.Port)
		log.Fatalf("PBX connection error [%s]: %s", amiConn, message)
	})

	err := a.ami.RegisterHandler(celTypeName, a.CELHandler)
	if err != nil {
		log.Error("AMI could not register handler: ", err)
	}

}

func (a *MyApp) CELHandler(m map[string]string) {
	log.Printf("CEL EVENT Received: %v\n", m)
	fields, err := getFields(m, celFieldEventName)
	if err != nil {
		log.Error("Error in CEL handler: ", err)
		return
	}
	log.Debug("Event CEL ", fields[celFieldEventName], " received")

	switch fields[celFieldEventName] {
	case celEventChanStart:
		fields, err := getFields(m, celFieldCallerIDnum, celFieldContext, celFieldUniqueId, celFieldLinkedId)
		if err != nil {
			log.Error("Error in parsing CEL ", celEventChanStart, ": ", err)
			return
		}
		if fields[celFieldContext] == a.config.Asterisk.Context && fields[celFieldUniqueId] == fields[celFieldLinkedId] {
			msg := fmt.Sprintf(msgCallIncoming, fields[celFieldCallerIDnum])
			if a.config.Crm.Enable {
				callerName, err := a.getCrmName(fields[celFieldCallerIDnum])
				if err != nil {
					log.Error("Error in requesting CRM: ", err)
				} else {
					if callerName != "" {
						msg = fmt.Sprintf("%s\n"+msgCallPerson, msg, callerName)
					}
				}
			}
			a.sendTelegramMessage(a.config.Telegram.ChatId, msg)
		}
	}

}

func getFields(m map[string]string, fields ...string) (map[string]string, error) {
	values := make(map[string]string)
	for _, field := range fields {
		value, ok := m[field]
		if !ok {
			log.WithFields(log.Fields{
				"map": m,
			}).Error("Invalid params map")
			// TODO FIX error handling here
			return nil, nil
		}
		values[field] = value
	}
	return values, nil
}
