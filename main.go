package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"iphoneBot/log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/robfig/cron"
	tb "gopkg.in/tucnak/telebot.v3"

	"github.com/joho/godotenv"
)

var Url = ""
var Token = ""
var MyId int64 = 0
var bot *tb.Bot
var isFetching = false
var cronTask *cron.Cron
var FetchTime = "3s"
var Proxy = ""
var storeMap = map[string]Store{}

type Store struct {
	Name      string
	Number    string
	Models    []Model
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string
	StoreNum  string
	StartTime time.Time
	MessageID int
	ChatID    int
	Enable    bool
}

func LoadEnv() {
	log.Log.Println("loadEnv")
	pwd, _ := os.Getwd()
	env := pwd + "/env"
	err := godotenv.Load(env)
	if err != nil {
		log.Log.Fatalf("unable to load env file")
	}
	MyId, _ = strconv.ParseInt(os.Getenv("MY_ID"), 10, 64)
	Url = os.Getenv("URL")
	Token = os.Getenv("BOT_TOKEN")
	FetchTime = os.Getenv("FETCH_DURATION")
	Proxy = os.Getenv("PROXY")
}

func main() {
	log.Init()
	LoadEnv()
	b, err := tb.NewBot(tb.Settings{
		Token:   Token,
		Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
		Verbose: false})

	if err != nil {
		fmt.Println("error " + err.Error())
		log.Log.Fatal(err)
		return
	}

	b.Handle("/ping", func(context tb.Context) error {
		return context.Send("pong!")
	})

	//调节查询时间  /fast 1s ==> 1秒查询一次
	b.Handle("/fast", func(context tb.Context) error {
		time := strings.TrimSpace(context.Message().Payload)
		err := startCron(time)
		if err != nil {
			return context.Reply(err.Error())
		}
		return nil
	})

	bot = b
	fetch()
	startCron(FetchTime)
	b.Start()
}

func startCron(time string) error {
	formatTime := fmt.Sprintf("@every %s", time)
	_, err := cron.Parse(formatTime)
	if err != nil {
		return err
	}
	if cronTask != nil {
		cronTask.Stop()
	}
	cronTask = cron.New()
	log.Log.Println(fmt.Sprintf("start cron time %s", formatTime))
	cronTask.AddFunc(formatTime, func() {
		fetch()
	})
	cronTask.Start()
	return nil
}

func fetch() {
	if isFetching {
		return
	}
	isFetching = true
	log.Log.Println(fmt.Sprintf("start fetch %v", time.Now().Local()))
	h := http.DefaultClient

	if len(Proxy) > 0 {
		proxyUrl, err := url.Parse(Proxy)
		if err != nil {
			log.Log.Error(err)
			return
		}
		h.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	h.Timeout = 2 * time.Second
	req, err := http.NewRequest(http.MethodGet, Url, nil)
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
		bot.Send(&tb.Chat{ID: MyId}, fmt.Sprintf("苹果获取状态不正常 %d %s", r.StatusCode, string(body)))
		return
	}

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		storeNum, _, _, err := jsonparser.Get(value, "storeNumber")
		if err != nil {
			log.Log.Println("error store num")
		}
		storeNumStr := string(storeNum)
		//过滤某个店
		// 		if storeNumStr == "R484" {
		// 			return
		// 		}
		//log.Log.Println(string(storeNum))
		storeName, err := jsonparser.GetString(value, "storeName")
		if err != nil {
			log.Log.Println(err)
		}
		storeName = strings.TrimSpace(storeName)

		partsAvailability, _, _, err := jsonparser.Get(value, "partsAvailability")
		if err != nil {
			return
		}
		//log.Log.Println(len(partsAvailability))
		//log.Log.Println(string(partsAvailability))
		if err != nil || len(partsAvailability) == 2 || len(partsAvailability) == 0 {
			clearStore(storeNumStr)
			isFetching = false
			return
		}
		hasModelList := []Model{}
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
			model := Model{Title: title, StoreNum: storeNumStr, StartTime: time.Now().Local(), Enable: true}
			hasModelList = append(hasModelList, model)
			return nil
		}, "partsAvailability")

		if len(hasModelList) == 0 {
			clearStore(storeNumStr)
			isFetching = false
			return
		}

		store, err := findStore(storeNumStr)

		cacheList := []Model{}

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
			store = Store{}
			store.Name = storeName
			store.Number = storeNumStr
			storeMap[storeNumStr] = store
		}
		for _, model := range hasModelList {
			if contains(cacheList, model) {
				continue
			}
			message, _ := notifyChannel(false, model)
			if message != nil {
				model.MessageID = message.ID
			}
			cacheList = append(cacheList, model)
		}
		store.Models = cacheList
		storeMap[storeNumStr] = store

		//log.Log.Println(cacheList)
		//log.Log.Printf("current availiable size = %d", len(cacheList))

	}, "body", "PickupMessage", "stores")

	isFetching = false
}

func notifyChannel(isDelete bool, model Model) (*tb.Message, error) {
	text := ""
	hasPreMessage := model.MessageID != -1
	if isDelete {
		text = getUnavailableStr(model, hasPreMessage)
	} else {
		storeName := model.StoreNum
		store, err := findStore(model.StoreNum)
		if err == nil {
			storeName = store.Name
		}
		text = getClickableStr(fmt.Sprintf("(%s)%s 点击购买", storeName, model.Title))
	}
	log.Log.Println(text)
	var err error
	if isDelete && hasPreMessage {
		_, err = bot.Reply(&tb.Message{ID: model.MessageID}, text, &tb.SendOptions{
			ParseMode:             tb.ModeMarkdown,
			DisableWebPagePreview: true,
		})
	} else {
		chatID := MyId
		message, err := bot.Send(&tb.Chat{ID: chatID}, text, &tb.SendOptions{
			ParseMode:             tb.ModeMarkdown,
			DisableWebPagePreview: true,
		})
		return message, err
	}
	if err != nil {
		log.Log.Error(err)
	}
	return nil, err
}

func deleteModel(list []Model, model Model) []Model {
	for i, m := range list {
		if m.Title == model.Title && m.Enable == model.Enable {
			notifyChannel(true, model)
			store, err := findStore(model.StoreNum)
			//从store中删除
			if err == nil {
				store.Models = deleteFromList(store.Models, model)
			}
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func deleteFromList(list []Model, model Model) []Model {
	for i, m := range list {
		if m.Title == model.Title {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func findStore(storeNumStr string) (Store, error) {
	store, ok := storeMap[storeNumStr]
	if ok {
		return store, nil
	}
	return Store{}, errors.New("no store")
}

func clearStore(storeNumStr string) bool {
	cacheStore, err := findStore(storeNumStr)
	if err != nil {
		return false
	}
	cacheModels := cacheStore.Models
	if len(cacheModels) > 0 {
		for _, model := range cacheModels {
			if model.Enable {
				notifyChannel(true, model)
				model.Enable = false
			}
		}
		cacheStore.Models = []Model{}
		log.Log.Println(cacheModels)
	}
	return false
}

func getUnavailableStr(model Model, hasPreMessage bool) string {
	nowTime := time.Now().Local()
	startTime := model.StartTime.Local()
	keepTime := nowTime.UnixMilli() - startTime.UnixMilli()
	timeStr := "毫秒"
	if keepTime > 1000 {
		keepTime = keepTime / 1000
		timeStr = "秒"
		if keepTime >= 60 {
			keepTime = keepTime / 60
			timeStr = "分钟"
			if keepTime >= 60 {
				keepTime = keepTime / 60
				timeStr = "小时"
			}
		}
	}
	var pre = ""
	if hasPreMessage {
		store, err := findStore(model.StoreNum)
		if err != nil {
			pre = fmt.Sprintf("(%s)%s", store.Name, model.Title)
		} else {
			pre = model.Title
		}
		pre += " "
	}
	return fmt.Sprintf("%s已被别人抢走，持续时间 %d %s，再接再厉", pre, keepTime, timeStr)
}

func getClickableStr(content string) string {
	return fmt.Sprintf("[%s](%s)", content, "https://www.apple.com.cn/shop/buy-iphone/iphone-13-pro/MLTE3CH/A")
}

func contains(list []Model, model Model) bool {
	for _, a := range list {
		if a.Title == model.Title && a.Enable == model.Enable {
			return true
		}
	}
	return false
}
