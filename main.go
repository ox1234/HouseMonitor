package main

import (
	"flag"
	"fmt"
	"houseMonitor/log"
	"houseMonitor/notify"
	"houseMonitor/notify/template"
	douban2 "houseMonitor/site/douban"
	"strings"
	"time"
)

func main() {
	var appID string
	var appSecret string
	var chatName string

	flag.StringVar(&appID, "a", "", "app id")
	flag.StringVar(&appSecret, "s", "", "app secret")
	flag.StringVar(&chatName, "c", "", "chat name")

	flag.Parse()

	if appID == "" || appSecret == "" || chatName == "" {
		flag.Usage()
		return
	}

	lk := notify.NewLarkNotify(appID, appSecret, chatName)
	douban := douban2.NewDouBanCollector()

	log.Info("start house monitor....")
	for range time.Tick(time.Minute * 10) {
		newItems, err := douban.Visit("https://www.douban.com/group/zhufang/discussion?start=0&type=new")
		if err != nil {
			log.Error("get house info fail: %s", err)
			continue
		}

		log.Info("%s get %d house item", time.Now(), len(newItems))
		if len(newItems) > 0 {
			var msgs []string
			for _, item := range newItems {
				msgs = append(msgs, item.String())
			}

			tpl := fmt.Sprintf(template.HouseTemplate, strings.Join(msgs, "\\n"))
			err = lk.SendMessage(tpl)
			if err != nil {
				log.Error("send message fail: %s", err)
			}
		}
	}

}
