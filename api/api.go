package api

import (
	"groupfibot/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReceiveMsgParm struct {
	Account string `json:"account"`
	GroupId string `json:"groupId"`
	Message string `json:"message"`
}

func RecieveMsg(c *gin.Context) {
	p := ReceiveMsgParm{}
	err := c.BindJSON(&p)
	if err != nil {
		slog.Error("recieve msg param error", "err", err)
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": 1,
			"err-msg":  "params error.",
		})
		return
	}

	if !service.ReceiveNewMsg(p.Account, p.GroupId, p.Message) {
		slog.Error("group is not exit", "account", p.Account, "groupid", p.GroupId)
		c.JSON(http.StatusOK, gin.H{
			"result":   false,
			"err-code": 2,
			"err-msg":  "group is not exit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": true,
	})
}
