package service

import (
	"crypto/rand"
	"fmt"
	"groupfibot/config"
	"log/slog"
	"math/big"
)

func SingleRandomResponse(msgs []string) (string, error) {
	if len(msgs) != 1 {
		return "", fmt.Errorf("msgs length isn't 1")
	}

	// 50%概率不回复；15%，进行10个单词以内的回复，20%，进行10-30个单词以内的回复；15%进行30-60个单词的回复
	num, _ := rand.Int(rand.Reader, big.NewInt(100))
	n := num.Int64()
	wn := int64(0)
	if n < 50 {
		wn = 0
	} else if n < 65 {
		wc, _ := rand.Int(rand.Reader, big.NewInt(9))
		wn = wc.Int64() + 1

	} else if n < 85 {
		wc, _ := rand.Int(rand.Reader, big.NewInt(20))
		wn = wc.Int64() + 10
	} else {
		wc, _ := rand.Int(rand.Reader, big.NewInt(30))
		wn = wc.Int64() + 30
	}
	extra := fmt.Sprintf(" limit %d words", wn)
	return GetBotResponse(msgs[0] + extra)
}

func LongTimeWaitResponse(msgs []string) (string, error) {
	if len(msgs) == 0 {
		return "", fmt.Errorf("msgs length is 0")
	}

	strs := []string{
		"hello", "good morning", "how's everything going?",
	}
	strs = append(strs, msgs[0]+" to the moon!")

	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(strs))))
	n := int(num.Uint64())
	return strs[n], nil
}

type SendMsgParam struct {
	Address string `json:"address"`
	GroupId string `json:"groupId"`
	Message string `json:"message"`
}

func SendMsg(account, groupId, msg string) bool {
	sm := SendMsgParam{
		Address: account,
		GroupId: groupId,
		Message: msg,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := HttpPost(config.SendUrl+"/send-message-to-group", sm, headers)
	if err != nil {
		slog.Error("Send message to group", "err", err, "param", sm)
		return false
	}
	slog.Info(string(res))
	return true
}

type BootstrapParam struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKeyHex"`
	GroupId    string `json:"groupId"`
}

func Bootstrap(account, groupid, pk string) bool {
	bm := BootstrapParam{
		Address:    account,
		GroupId:    groupid,
		PrivateKey: pk,
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	res, err := HttpPost(config.SendUrl+"/bootstrap-and-enter-group", bm, headers)
	if err != nil {
		slog.Error("bootstrap a group", "err", err, "param", bm)
		return false
	}
	slog.Info(string(res))
	return true
}
