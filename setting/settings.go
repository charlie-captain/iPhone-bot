package setting

import (
	"encoding/json"
	"fmt"
	"iphoneBot/log"
	"iphoneBot/model"
	"net/http"
	"net/url"
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
	Models        []string           `json:"model_list"`
	Region        string             `json:"region"`
}

var Host = "https://www.apple.com"

const RecommendUrlSuffix = "/shop/pickup-message-recommendations?t=compact&searchNearby=true&store=%s&product="
const IPhoneModelSuffix = "/shop/fulfillment-messages?store=%s&little=false&parts.0=%s&mts.0=regular&mts.1=sticky&fts=true&searchNearby=true"
const IPhoneBuySuffix = "/shop/buy-iphone/"
const CurIPhone = "iphone-15-pro"

const AUTO_DELETE_TIME = 10 * time.Second

// iPhone15 Pro
var IPhoneProModelList = []string{
	"MTQJ3CH/A", //1T pro
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
	//log.Log.Println(settings)
	if len(settings.Stores) == 0 {
		log.Log.Fatal("store must not be empty")
	}

	if settings.Region == "" {
		//默认中国
		Host = Host + ".cn"
	} else if len(settings.Region) > 0 {
		Host = Host + "/" + settings.Region
	}
	fetchSuffix := fmt.Sprintf(RecommendUrlSuffix, settings.Stores[0])
	exactlyMode := len(settings.Models) > 0
	if exactlyMode {
		fetchSuffix = IPhoneModelSuffix
		IPhoneProModelList = settings.Models
	}
	fetchSource := &model.FetchSource{
		Url:         Host + fetchSuffix,
		Type:        IPhoneProModelList,
		ExactlyMode: exactlyMode,
	}
	settings.FetchSource = fetchSource
	return settings
}

func SetUpProxy(settings *Settings, h *http.Client) bool {
	if len(settings.Proxy) > 0 {
		proxyUrl, err := url.Parse(settings.Proxy)
		if err != nil {
			log.Log.Error(err)
			return true
		}
		h.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	return false
}
