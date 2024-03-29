package main

const (
	configFileName = "conf.yml"
	configLogDebug = "debug"
	configLogInfo  = "info"
	configLogWarn  = "warn"
	configLogError = "error"
	configLogFatal = "fatal"
)

const (
	celTypeName         = "CEL"
	celEventChanStart   = "CHAN_START"
	celFieldEventName   = "EventName"
	celFieldCallerIDnum = "CallerIDnum"
	celFieldContext     = "Context"
	celFieldUniqueId    = "UniqueID"
	celFieldLinkedId    = "LinkedID"
)
const (
	imapFolder   = "INBOX"
	imapFlagSeen = "\\Seen"
)

const (
	tgCommandStatus = "status"
	tgCommandTimer  = "timer"
	tgCommandWho    = "who"
)

const (
	timerSeconds            = "s"
	timerSecondsDescription = "сек"
	timerMinutes            = "m"
	timerMinutesDescription = "мин"
	timerHours              = "h"
	timerHoursDescription   = "час"
)

//TODO добавить смайлов
const (
	msgCallIncoming    = "\u260e Входящий звонок с номера %s"
	msgCallPerson      = "\xF0\x9F\x91\xA4 %s"
	msgMailIncoming    = "\xF0\x9F\x93\xA8 %s"
	msgStatusConnected = "\u270c Привет! Я в сети ;)"
	msgStatusAlive     = "\xF0\x9F\x99\x8B Я в порядке!"
	msgUnknownCommand  = "Не понимаю, что ты от меня хочешь \xF0\x9F\x99\x87"
	msgTimerStarted    = "\xF0\x9F\x91\x8C Напомню через %d %s"
	msgTimerFinished   = "\xE2\x8F\xB0 Напоминалка. Время вышло!"
	msgTimerFailed     = "\xE2\x84\xB9 Команду /timer надо использовать хотя бы с 1 аргументом.\n " +
		"Первый аргумент - это время, через которое нужно напомнить.\n" +
		"Второй аргумент (не обязательно) - это в каких единицах измеряется время: " + timerSeconds + " (секунды), " + timerMinutes + " (минуты, по-умолчанию), " + timerHours + " (часы). \n" +
		"Например, чтобы напомнить через 5 минут, надо написать\n" +
		"`/timer 5 " + timerMinutes + "`"
	msgWhoNoCRM  = "\xF0\x9F\x99\x87 Интеграция с CRM не активна."
	msgWhoReply  = "\xF0\x9F\x94\x8E Поиск в CRM номера %s"
	msgWhoFailed = "\xE2\x84\xB9 Команду /who надо использовать с одним аргументом.\n" +
		"Это аргумент - номер телефона, по которому осуществляется поиск в CRM\n" +
		"Например, чтобы найти кто звонил с номера +74957485588, надо написать\n" +
		"`/who +74957485588`"
)
