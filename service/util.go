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

func clearStore(storeNumStr string, modelType string) bool {
	cacheStore, err := findStore(storeNumStr)
	if err != nil {
		return false
	}
	cacheModels := cacheStore.Models
	if len(cacheModels) > 0 {
		for i, model := range cacheModels {
			if model.Enable && model.ModelName == modelType {
				bot.NotifyChannel(true, _settings.ChatID, model, cacheStore)
				model.Enable = false
				cacheStore.Models = append(cacheStore.Models[:i], cacheStore.Models[i+1:]...)
			}
		}
		storeMap[storeNumStr] = cacheStore
		log.Log.Println(cacheModels)
	}
	return false
}
