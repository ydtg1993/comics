package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/tools"
	"comics/tools/config"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tidwall/gjson"
	"os"
	"strconv"
	"strings"
	"time"
)

func ComicPaw() {
	tags := map[string]int{
		"恋爱": 20,
		"古风": 46,
		"穿越": 80,
	}
	regions := map[string]int{
		"国漫": 2,
		"韩漫": 3,
		"日漫": 4,
	}
	pays := map[string]int{
		"免费": 1,
		"付费": 2,
	}
	states := map[string]int{
		"连载中": 1,
		"已完结": 2,
	}
	for tag, tagId := range tags {
		for region, regionId := range regions {
			for pay, payId := range pays {
				for state, stateId := range states {
					fmt.Printf("%s %s %s %s \n", tag, region, pay, state)
					category(tagId, regionId, payId, stateId)
				}
			}
		}
	}
}

func ComicUpdate() {

}

func category(tagId, regionId, payId, stateId int) {
	url := fmt.Sprintf(os.Getenv("SOURCE_URL")+"/tag/%d?region=%d&pays=%d&state=%d&sort=1&page=1",
		tagId, regionId, payId, stateId)
	header := map[string]string{
		"User-Agent": config.Spe.UserAgent,
		"Referer":    config.Spe.SourceUrl}
	content, _, code := tools.HttpRequestByHeaderFor5(url, "GET", "", header)
	if code != 200 {
		logs.Error("无法抓取目标页:" + url)
		return
	}

	t := time.NewTicker(time.Second * time.Duration(5))
	<-t.C
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		logs.Error("无法抓取目标页Dom:" + url)
		return
	}

	lastPage := doc.Find("ul.pagination>.itemBten").Last().Text()
	if lastPage == "" {
		paw(tagId, regionId, payId, stateId, 1)
		return
	}
	last, _ := strconv.Atoi(strings.TrimSpace(lastPage))
	for ; last > 0; last-- {
		paw(tagId, regionId, payId, stateId, last)
	}
}

func paw(tagId, regionId, payId, stateId, last int) {
	header := map[string]string{
		"User-Agent": config.Spe.UserAgent,
		"Referer":    config.Spe.SourceUrl}
	url := fmt.Sprintf(os.Getenv("SOURCE_URL")+"/search/mini/topic/multi_filter?tag_id=%d&label_dimension_origin=%d&pay_status=%d&update_status=%d&sort=1&page=%d&size=48",
		tagId, regionId, payId, stateId, last)
	content, _, code := tools.HttpRequestByHeaderFor5(url, "GET", "", header)
	if code != 200 {
		logs.Error("无法抓取目标页 分页:" + url)
		return
	}

	t := time.NewTicker(time.Second * time.Duration(5))
	<-t.C

	list := gjson.Get(content, "hits.topicMessageList")
	list.ForEach(func(key, value gjson.Result) bool {
		id, _ := strconv.Atoi(value.Get("id").String())

		sourceComic := new(model.SourceComic)
		sourceComic.Source = 1
		sourceComic.SourceId = id
		sourceComic.Cover = value.Get("cover_image_url").String()
		sourceComic.SourceUri = "web/topic/" + value.Get("id").String()
		sourceComic.Title = value.Get("title").String()

		orm.Eloquent.Where("source = ? and source_id = ?", 1, id).FirstOrCreate(&sourceComic)
		return true
	})
}
