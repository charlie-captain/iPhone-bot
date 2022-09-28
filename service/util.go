package service

import (
	"errors"
	"iphoneBot/bot"
	"iphoneBot/log"
	"iphoneBot/model"
)

func findStore(storeNumStr string) (model.Store, error) {
	store, ok := storeMap[storeNumStr]
	if ok {
		return store, nil
	}
	return model.Store{}, errors.New("no store")
}

func clearStore(storeNumStr string, modelType []string) bool {
	cacheStore, err := findStore(storeNumStr)
	if err != nil {
		return false
	}
	cacheModels := cacheStore.Models
	if len(cacheModels) > 0 {
		for _, model := range cacheModels {
			if model.Enable {
				bot.NotifyChannel(true, _settings.ChatID, model, cacheStore)
				model.Enable = false
			}
		}
		cacheStore.Models = []model.Model{}
		log.Log.Println(cacheModels)
	}
	return false
}
