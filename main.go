package main

import (
	"fmt"
	"houseMonitor/log"
	"houseMonitor/notify"
	"houseMonitor/notify/template"
	douban2 "houseMonitor/site/douban"
	"strings"
	"time"
)

func main() {
	lk := notify.NewLarkNotify("cli_a2649b64eb79900e", "WDDkcui55FJgrF6TL1yfQhVpL2vS1KgT", "北京租房监控")
	douban := douban2.NewDouBanCollector()

	for range time.Tick(time.Second * 30) {
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
