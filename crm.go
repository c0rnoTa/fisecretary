package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

func (a *MyApp) getCrmName(phone string) (string, error) {
	log.Info("CRM request caller name for number: ", phone)
	client := http.Client{
		Timeout: time.Duration(a.config.Crm.Timeout) * time.Second,
	}
	log.Debug("CRM GET ", fmt.Sprintf(a.config.Crm.Url, phone))
	resp, err := client.Get(fmt.Sprintf(a.config.Crm.Url, phone))
	if err != nil {
		log.Error("CRM request error: ", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	log.Debug("Raw CRM response", string(bodyData))
	if err != nil {
		log.Error("CRM can't get response", err)
		return "", err
	}

	return string(bodyData), nil
}
