package basic

import (
	"fmt"
	"time"
)

type HouseMetaData struct {
	HouseURL       string
	Description    string
	SubmitPerson   string
	SubmitPersonID string
	PostTime       time.Time
}

func (h *HouseMetaData) String() string {
	return fmt.Sprintf("🏠 [%s](%s)  发布日期：%s  发布人： %s", h.Description, h.HouseURL, h.PostTime.Format("01-02 15:04"), h.SubmitPerson)
}

type HouseCache struct {
	Cache []*HouseMetaData
}

func (h *HouseCache) Add(meta *HouseMetaData) bool {
	for _, item := range h.Cache {
		if item.HouseURL == meta.HouseURL {
			return false
		}
	}
	h.Cache = append(h.Cache, meta)
	return true
}
