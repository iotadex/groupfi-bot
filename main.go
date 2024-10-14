package main

import (
	"groupfibot/api"
	"groupfibot/config"
	"groupfibot/daemon"
	"groupfibot/service"
	"log/slog"
	"os"

	"github.com/wytools/rlog/handler"
)

func main() {
	if os.Args[len(os.Args)-1] != "-d" {
		os.Args = append(os.Args, "-d")
	}
	daemon.Background("./out.log", true)

	config.Load()

	slog.SetDefault(handler.GetDefaultDailyLogger("logs/out.log", 0, 0))

	api.StartHttpServer(config.HttpPort)

	service.StartAllChats()

	daemon.WaitForKill()
}
