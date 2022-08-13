package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"houseMonitor/log"
	"houseMonitor/notify"
	"houseMonitor/notify/template"
	douban2 "houseMonitor/site/douban"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	var appID string
	var appSecret string
	var chatName string
	var conf string

	flag.StringVar(&appID, "a", "", "app id")
	flag.StringVar(&appSecret, "s", "", "app secret")
	flag.StringVar(&chatName, "c", "", "chat name")
	flag.StringVar(&conf, "conf", "", "config file")

	flag.Parse()

	if appID == "" || appSecret == "" || chatName == "" {
		flag.Usage()
		return
	}

	lk := notify.NewLarkNotify(appID, appSecret, chatName)
	douban := douban2.NewDouBanCollector()

	banWords, err := getBanList(conf)
	if err != nil {
		log.Error("get ban list fail: %s", err)
		return
	}

	go func() {
		for range time.Tick(time.Second * 30) {
			log.Info("start get ban list...")
			banWords, err = getBanList(conf)
			if err != nil {
				log.Error("cron get ban list fail: %s", err)
				continue
			}
			log.Info("refresh ban list with: %+v", banWords)
		}
	}()

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
				for _, banWord := range banWords {
					if strings.Contains(item.Description, banWord) {
						continue
					}
				}
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

func getBanList(u string) ([]string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("get ban list fail: %w", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get body fail: %w", err)
	}
	var banList map[string][]string
	err = json.Unmarshal(b, &banList)
	if err != nil {
		return nil, fmt.Errorf("unmarshal body fail: %w", err)
	}
	return banList["ban_words"], nil
}
