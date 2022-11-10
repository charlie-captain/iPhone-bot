package bot

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"iphoneBot/log"
	"iphoneBot/model"
	"iphoneBot/setting"
	"net/http"
	"time"
)

var bot *tb.Bot

func Init(settings *setting.Settings) *tb.Bot {
	client := http.DefaultClient
	client.Timeout = 10 * time.Second
	setting.SetUpProxy(settings, client)
	b, err := tb.NewBot(tb.Settings{
		Token:   settings.BotToken,
		Poller:  &tb.LongPoller{Timeout: 10 * time.Second},
		Verbose: false,
		Client:  client})

	if err != nil {
		fmt.Println("error " + err.Error())
		log.Log.Fatal(err)
		return nil
	}

	b.Handle("/ping", func(context tb.Context) error {
		return context.Send("pong!")
	})

	bot = b
	return bot
}

func NotifyChannel(isDelete bool, chatID int64, newModel model.Model, store model.Store) (*tb.Message, error) {
	text := ""
	hasPreMessage := newModel.MessageID != -1
	if isDelete {
		text = getUnavailableStr(newModel, store)
	} else {
		storeName := store.Name
		text = fmt.Sprintf("🎉（%s）%s 点击购买", storeName, newModel.Title)
	}
	modelName := newModel.ModelName
	text = getClickableStr(modelName, text)
	log.Log.Println(text)
	var err error
	if isDelete && hasPreMessage {
		replyMsg, err := bot.Reply(&tb.Message{ID: newModel.MessageID, Chat: &tb.Chat{ID: chatID}}, text, &tb.SendOptions{
			ParseMode:             tb.ModeMarkdown,
			DisableWebPagePreview: true,
		})
		if err != nil {
			log.Log.Error(err)
		} else {
			run(bot, DelayTask{
				chatID:        chatID,
				editMessage:   newModel.MessageID,
				deleteMessage: replyMsg.ID,
				text:          text,
				delayTime:     setting.AUTO_DELETE_TIME})
		}
	} else {
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

func getUnavailableStr(model model.Model, store model.Store) string {
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
	storeName := ""
	storeName = store.Name
	pre = fmt.Sprintf("（%s）%s", storeName, model.Title)
	pre += " "
	return fmt.Sprintf("💪🏻 %s已被抢走，持续 %d %s，再接再厉", pre, keepTime, timeStr)
}

func getClickableStr(modelType string, content string) string {
	url := setting.Host + setting.IPhoneBuySuffix + setting.CurIPhone + "/" + modelType
	return fmt.Sprintf("[%s](%s)", content, url)
}

type DelayTask struct {
	chatID        int64
	editMessage   int
	deleteMessage int
	text          string
	startTime     int64
	delayTime     time.Duration
}

func run(bot *tb.Bot, delayTask DelayTask) {
	time.AfterFunc(delayTask.delayTime, func() {
		err := bot.Delete(&tb.Message{ID: delayTask.deleteMessage, Chat: &tb.Chat{ID: delayTask.chatID}})
		if err != nil {
			log.Log.Error("删除失败" + err.Error())
			// 删除另一条
			err := bot.Delete(&tb.Message{ID: delayTask.editMessage, Chat: &tb.Chat{ID: delayTask.chatID}})
			if err != nil {
				log.Log.Error("都删除失败" + err.Error())
			}
		} else {
			message, err := bot.Edit(&tb.Message{ID: delayTask.editMessage, Chat: &tb.Chat{ID: delayTask.chatID}},
				delayTask.text, &tb.SendOptions{ParseMode: tb.ModeMarkdown, DisableWebPagePreview: true})
			if err != nil {
				fmt.Println("修改失败 " + err.Error())
			} else {
				messageId, chatId := message.MessageSig()
				fmt.Println("修改成功, id = " + messageId + " chatId = " + fmt.Sprintf("%d", chatId))
			}
		}
	})
}
