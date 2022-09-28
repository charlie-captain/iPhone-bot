package service

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/robfig/cron"
	"io/ioutil"
	"iphoneBot/bot"
	"iphoneBot/log"
	"iphoneBot/model"
	"iphoneBot/setting"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var isFetching = false
var cronTask *cron.Cron
var storeMap = map[string]model.Store{}
var _settings *setting.Settings

func Init(settings *setting.Settings) {
	_settings = settings
}

func StartCron(time string) error {
	formatTime := fmt.Sprintf("@every %s", time)
	_, err := cron.Parse(formatTime)
	if err != nil {
		log.Log.Error(err)
		return err
	}
	if cronTask != nil {
		cronTask.Stop()
	}
	cronTask = cron.New()
	log.Log.Println(fmt.Sprintf("start cron time %s", formatTime))
	cronTask.AddFunc(formatTime, func() {
		Fetch(_settings.FetchSource)
	})
	cronTask.Start()
	return nil
}

func Fetch(source *model.FetchSource) {
	if isFetching {
		return
	}
	isFetching = true
	log.Log.Println(fmt.Sprintf("start fetch %v", time.Now().Local()))
	h := http.DefaultClient

	if len(_settings.Proxy) > 0 {
		proxyUrl, err := url.Parse(_settings.Proxy)
		if err != nil {
			log.Log.Error(err)
			return
		}
		h.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	h.Timeout = 2 * time.Second
	url := source.Url + source.Type[0]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Log.Println(err)
		isFetching = false
		return
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.49 Safari/537.36")
	req.Header.Set("authority", "www.apple.com.cn")
	h.Timeout = 10 * time.Second

	r, err := h.Do(req)
	if r != nil {
		defer r.Body.Close()
	}
	if err != nil {
		log.Log.Println(err)
		isFetching = false
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Log.Println(err)
		isFetching = false
		return
	}
	log.Log.Println(fmt.Sprintf("status %s, length %d", r.Status, r.ContentLength))

	if r.StatusCode != 200 {
		//不正常
		return
	}

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		storeNum, _, _, err := jsonparser.Get(value, "storeNumber")
		if err != nil {
			log.Log.Println("error store num")
		}
		storeNumStr := string(storeNum)
		//过滤店铺
		var hasStore = false
		for _, settingStore := range _settings.Stores {
			if storeNumStr == settingStore {
				hasStore = true
				break
			}
		}
		if !hasStore {
			return
		}
		storeName, err := jsonparser.GetString(value, "storeName")
		if err != nil {
			log.Log.Error(err)
		}
		storeName = strings.TrimSpace(storeName)

		partsAvailability, _, _, err := jsonparser.Get(value, "partsAvailability")
		if err != nil {
			log.Log.Error(err)
			return
		}
		if err != nil || len(partsAvailability) == 2 || len(partsAvailability) == 0 {
			clearStore(storeNumStr, source.Type)
			isFetching = false
			return
		}
		hasModelList := []model.Model{}
		jsonparser.ObjectEach(value, func(key, value []byte, dataType jsonparser.ValueType, offset int) error {
			messageValue, _, _, err := jsonparser.Get(value, "messageTypes", "regular")
			if err != nil {
				log.Log.Error(err)
				return err
			}
			title, err := jsonparser.GetUnsafeString(messageValue, "storePickupProductTitle")
			if err != nil {
				log.Log.Println(err)
				return err
			}
			available, err := jsonparser.GetUnsafeString(value, "pickupDisplay")
			if err != nil {
				log.Log.Println(err)
				return err
			}
			if string(available) == "unavailable" {
				log.Log.Println(fmt.Sprintf("%s 不可取货", title))
				return nil
			}
			modelName, err := jsonparser.GetUnsafeString(value, "partNumber")
			if err != nil {
				log.Log.Error(err)
				return err
			}
			model := model.Model{Title: title, StoreNum: storeNumStr, StartTime: time.Now().Local(), Enable: true, ModelName: modelName}
			hasModelList = append(hasModelList, model)
			return nil
		}, "partsAvailability")

		if len(hasModelList) == 0 {
			clearStore(storeNumStr, nil)
			isFetching = false
			return
		}

		store, err := findStore(storeNumStr)

		cacheList := []model.Model{}

		if err == nil {
			//存在list, 进行去重
			cacheList = store.Models
			for _, model := range cacheList {
				if contains(hasModelList, model) {
					continue
				}
				//删除不包含的
				cacheList = deleteModel(cacheList, model)
			}
		} else {
			store = model.Store{}
			store.Name = storeName
			store.Number = storeNumStr
			storeMap[storeNumStr] = store
		}
		for _, model := range hasModelList {
			if contains(cacheList, model) {
				continue
			}
			store, err := findStore(model.StoreNum)
			if err != nil {
				log.Log.Error(err)
			}
			message, _ := bot.NotifyChannel(false, _settings.ChatID, model, store)
			if message != nil {
				model.MessageID = message.ID
			}
			cacheList = append(cacheList, model)
		}
		store.Models = cacheList
		storeMap[storeNumStr] = store
	}, "body", "PickupMessage", "stores")

	isFetching = false
}

func deleteModel(list []model.Model, model model.Model) []model.Model {
	for i, m := range list {
		if m.Title == model.Title && m.Enable == model.Enable {
			store, err := findStore(model.StoreNum)
			//从store中删除
			if err == nil {
				store.Models = deleteFromList(store.Models, model)
			}
			if m.Enable {
				bot.NotifyChannel(true, _settings.ChatID, model, store)
			}
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func deleteFromList(list []model.Model, model model.Model) []model.Model {
	for i, m := range list {
		if m.Title == model.Title {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func contains(list []model.Model, model model.Model) bool {
	for _, a := range list {
		if a.Title == model.Title && a.Enable == model.Enable {
			return true
		}
	}
	return false
}
