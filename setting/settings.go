package setting

import (
	"encoding/json"
	"fmt"
	"iphoneBot/log"
	"iphoneBot/model"
	"os"
	"time"
)

type Settings struct {
	BotToken      string             `json:"bot_token"`
	ChatID        int64              `json:"chat_id"`
	Stores        []string           `json:"store_list"`
	FetchDuration string             `json:"fetch_duration"`
	Proxy         string             `json:"proxy"`
	FetchSource   *model.FetchSource `json:"-"`
}

const IPhoneProUrl = "https://www.apple.com.cn/shop/pickup-message-recommendations?t=compact&searchNearby=true&store=%s&product="
const IPhoneProBuyUrl = "https://www.apple.com.cn/shop/buy-iphone/"
const CurIPhone = "iphone-14-pro"

const AUTO_DELETE_TIME = 10 * time.Second

// iPhone14 Pro 灰色
var IPhoneProModelList = []string{
	"MQ1C3CH/A", //256G
}

func LoadEnv() *Settings {
	log.Log.Println("loadEnv")
	pwd, _ := os.Getwd()
	settingsFile := pwd + "/settings.json"
	file, err := os.Open(settingsFile)
	defer file.Close()
	if err != nil {
		log.Log.Fatal(err)
		return nil
	}
	var settings *Settings
	jsDecoder := json.NewDecoder(file)
	err = jsDecoder.Decode(&settings)
	if err != nil {
		log.Log.Fatalf("unable to load settings.json file")
	}
	log.Log.Println(settings)
	settings.FetchSource = &model.FetchSource{
		Url:  fmt.Sprintf(IPhoneProUrl, settings.Stores[0]),
		Type: IPhoneProModelList,
	}
	return settings
}
