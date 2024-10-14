package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	HttpPort     int           // http service port
	SendUrl      string        // send msg url
	ChatAccounts []ChatAccount // chat accounts
)

type ChatAccount struct {
	Account    string `json:"account"`
	PrivateKey string `json:"private_key"`
	GroupId    string `json:"groupid"`
	TokenName  string `json:"token_name"`
	MsgNumber  int    `json:"msg_number"`
	Strategy   int    `json:"strategy"`
}

// Load load config file
func Load() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	type Config struct {
		HttpPort     int           `json:"http_port"`
		SendUrl      string        `json:"send_url"`
		ChatAccounts []ChatAccount `json:"chats"`
	}
	all := &Config{}
	if err = json.NewDecoder(file).Decode(all); err != nil {
		log.Panic(err)
	}
	HttpPort = all.HttpPort
	SendUrl = all.SendUrl
	ChatAccounts = all.ChatAccounts
}
