package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Структура конфигурационного файла
type Config struct {
	Imap struct {
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		Server         string `yaml:"server"`
		RefreshTimeout int64  `yaml:"refresh"`
	} `yaml:"imap"`
	Telegram struct {
		Token  string `yaml:"token"`
		ChatId int64  `yaml:"chatid"`
	} `yaml:"telegram"`
	Asterisk struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Context  string `yaml:"context"`
	} `yaml:"asterisk"`
	Crm struct {
		Url     string `yaml:"url"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"crm"`
	LogLevel string `yaml:"loglevel"`
}

// Читаем конфиг и устанавливаем параметры приложения
func (a *MyApp) GetConfigYaml(filename string) {
	log.Info("Reading config ", filename)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, &a.config)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	a.logLevel = setLogLevel(a.config.LogLevel)
}

// Устанавливаем уровень журналирования событий в приложении
func setLogLevel(confLogLevel string) log.Level {
	var result log.Level
	switch confLogLevel {
	case configLogDebug:
		result = log.DebugLevel
	case configLogInfo:
		result = log.InfoLevel
	case configLogWarn:
		result = log.WarnLevel
	case configLogError:
		result = log.ErrorLevel
	case configLogFatal:
		result = log.FatalLevel
	default:
		result = log.InfoLevel
	}

	log.Info("Application logging level: ", result)

	return result
}
