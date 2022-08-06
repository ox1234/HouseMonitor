package douban

import (
	"fmt"
	"github.com/gocolly/colly"
	"houseMonitor/basic"
	"houseMonitor/log"
	"time"
)

type DoubanCollector struct {
	c      *colly.Collector
	cache  *basic.HouseCache
	NewAdd []*basic.HouseMetaData
}

func NewDouBanCollector() *DoubanCollector {
	douban := new(DoubanCollector)
	douban.cache = new(basic.HouseCache)

	c := colly.NewCollector()
	c.AllowURLRevisit = true
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		if e.Attr("class") == "" {
			ch := e.DOM.Children()
			if ch.Length() == 4 {
				houseTitle, _ := ch.Eq(0).Children().Eq(0).Attr("title")
				houseHref, _ := ch.Eq(0).Children().Eq(0).Attr("href")

				author := ch.Eq(1).Children().Eq(0).Text()
				if author == "草原" {
					return
				}
				authorHref, _ := ch.Eq(1).Children().Eq(0).Attr("href")

				postTime := ch.Eq(3).Text()
				fixTime := fmt.Sprintf("2022-%s:00", postTime)
				t, err := time.ParseInLocation("2006-01-02 15:04:05", fixTime, time.Local)
				if err != nil {
					log.Error("parse %s time fail: %s", fixTime, err)
					return
				}

				meta := &basic.HouseMetaData{
					HouseURL:       houseHref,
					Description:    houseTitle,
					SubmitPerson:   author,
					SubmitPersonID: authorHref,
					PostTime:       t,
				}

				isNew := douban.cache.Add(meta)
				if isNew {
					douban.NewAdd = append(douban.NewAdd, meta)
				}
			}
		}
	})

	douban.c = c

	return douban
}

func (d *DoubanCollector) Visit(url string) ([]*basic.HouseMetaData, error) {
	d.NewAdd = []*basic.HouseMetaData{}
	err := d.c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("visit %s fail: %w", url, err)
	}
	return d.NewAdd, nil
}
