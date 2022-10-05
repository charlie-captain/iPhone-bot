package main

import (
	"iphoneBot/bot"
	"iphoneBot/log"
	"iphoneBot/service"
	"iphoneBot/setting"
)

func main() {
	log.Init()
	settings := setting.LoadEnv()
	b := bot.Init(settings)
	service.Init(settings)
	service.StartFetchAll()
	service.StartCron(settings.FetchDuration)
	b.Start()
}
