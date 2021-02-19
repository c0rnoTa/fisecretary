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
)

const (
	msgCallIncoming    = "Входящий звонок с номера %s"
	msgMailIncoming    = "%s"
	msgStatusConnected = "\u270c Привет! Я в сети ;)"
	msgStatusAlive     = "\u2705 Я в порядке!"
	msgUnknownCommand  = "Не понимаю, что ты от меня хочешь :("
)
