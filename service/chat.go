package service

import (
	"crypto/rand"
	"groupfibot/config"
	"log/slog"
	"math/big"
	"sync"
	"time"
)

var allChats map[string]*BotChat // key : account-groupid
func StartAllChats() {
	allChats = make(map[string]*BotChat)
	for _, c := range config.ChatAccounts {
		if Bootstrap(c.Account, c.GroupId, c.PrivateKey) {
			bc := NewBotChat(c.Account, c.PrivateKey, c.GroupId, c.MsgNumber)
			allChats[c.Account+"-"+c.GroupId] = bc
		}
	}
}

func ReceiveNewMsg(account, groupid, msg string) bool {
	key := account + "-" + groupid
	if bc, exist := allChats[key]; exist {
		bc.ReceiveNewMsg(msg)
		return true
	}
	return false
}

type BotChat struct {
	Account       string   // evm address
	PrivateKey    string   // account's private key
	Groupid       string   // group id
	TokenName     string   // gourp's token name
	LatestTime    int64    // the latest message's timestamp as seconds
	Messages      []string // the latest messages
	MessageNumber int      // the number of latest messages to store
	newMsgSignal  chan bool
	sync.Mutex
}

func NewBotChat(account, privateKey, groupid string, msgNumber int) *BotChat {
	bc := &BotChat{
		Account:       account,
		PrivateKey:    privateKey,
		Groupid:       groupid,
		LatestTime:    time.Now().Unix(),
		Messages:      make([]string, 0, msgNumber+1),
		MessageNumber: msgNumber,
		newMsgSignal:  make(chan bool, 5),
	}
	go bc.Run()
	return bc
}

func (bc *BotChat) ReceiveNewMsg(msg string) {
	bc.Lock()
	defer bc.Unlock()
	bc.Messages = append(bc.Messages, msg)
	if len(bc.Messages) > bc.MessageNumber {
		bc.Messages = bc.Messages[1:]
	}
	bc.LatestTime = time.Now().Unix()
	// perform the strategy
	bc.newMsgSignal <- true
}

func (bc *BotChat) Run() {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-bc.newMsgSignal:
			msg := bc.Messages[len(bc.Messages)-1]
			msg, err := SingleRandomResponse([]string{msg})
			if err != nil {
				slog.Error("Single random response", "err", err)
			} else {
				SendMsg(bc.Account, bc.Groupid, msg)
			}
		case <-ticker.C:
			bc.Lock()
			latestTs := bc.LatestTime
			bc.Unlock()
			nowTs := time.Now().Unix()
			if (nowTs - latestTs) > 43200 { // over 12 hours
				go func() {
					num, _ := rand.Int(rand.Reader, big.NewInt(1200))
					time.Sleep(time.Second * time.Duration(num.Int64()))
					msg, _ := LongTimeWaitResponse([]string{bc.TokenName})
					if SendMsg(bc.Account, bc.Groupid, msg) {
						bc.Lock()
						bc.LatestTime = time.Now().Unix()
						bc.Unlock()
					}
				}()
			}
		}
	}
}
