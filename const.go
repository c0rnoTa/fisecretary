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
)

const (
	timerSeconds            = "s"
	timerSecondsDescription = "секунд"
	timerMinutes            = "m"
	timerMinutesDescription = "минут"
	timerHours              = "h"
	timerHoursDescription   = "часов"
)

//TODO добавить смайлов
const (
	msgCallIncoming    = "\u260e Входящий звонок с номера %s"
	msgCallPerson      = "\xF0\x9F\x91\xA4 %s"
	msgMailIncoming    = "\xF0\x9F\x93\xA8 %s"
	msgStatusConnected = "\u270c Привет! Я в сети ;)"
	msgStatusAlive     = "\xF0\x9F\x99\x8B Я в порядке!"
	msgUnknownCommand  = "Не понимаю, что ты от меня хочешь \xF0\x9F\x99\x87"
	msgTimerStarted    = "Напомню через %d %s"
	msgTimerFinished   = "Время вышло!!! Напоминаю!!!"
	msgTimerFailed     = "Команду /timer надо использовать хотя бы с 1 аргументом. \n " +
		"Первый аргумент - это время, через которое нужно напомнить. \n" +
		"Второй аргумент (не обязательно) - это в каких единицах измеряется время: " + timerSeconds + " (секунды), " + timerMinutes + " (минуты), " + timerHours + " (часы). \n" +
		"Например, чтобы напомнить через 5 минут, надо написать `/timer 5 " + timerMinutes + "`"
)
