package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

func ComicPaw() {
	tags := map[string]int{
		"恋爱": 20,
		//"古风": 46,
		//"穿越": 80,
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
	bot := colly.NewCollector(
		colly.AllowedDomains(config.Spe.SourceUrl),
	)
	extensions.RandomUserAgent(bot)
	extensions.Referer(bot)

	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/tag/%d?region=%d&pays=%d&state=%d&sort=1&page=1",
		tagId, regionId, payId, stateId)
	lastPage := 1
	bot.OnHTML("ul.pagination", func(e *colly.HTMLElement) {
		page := e.DOM.Find(".itemBten").Last().Text()
		if page != "" {
			lastPage, _ = strconv.Atoi(strings.TrimSpace(page))
		}
		paw(tagId, regionId, payId, stateId, lastPage)
	})
	bot.OnResponse(func(r *colly.Response) {

	})
	err := bot.Visit(url)
	if err != nil {
		logs.Error("无法抓取分类列表页Dom:" + url)
	}
}

func paw(tagId, regionId, payId, stateId, last int) {
	url := fmt.Sprintf("https://"+config.Spe.SourceUrl+"/search/mini/topic/multi_filter?tag_id=%d&label_dimension_origin=%d&pay_status=%d&update_status=%d&sort=1&page=%d&size=48",
		tagId, regionId, payId, stateId, last)
	content, err := requestApi(url, "GET", "", 3)
	if err != nil {
		return
	}

	list := content.Get("hits.topicMessageList")
	list.ForEach(func(key, value gjson.Result) bool {
		id, _ := strconv.Atoi(value.Get("id").String())

		sourceComic := new(model.SourceComic)
		sourceComic.Source = 1
		sourceComic.SourceId = id
		sourceComic.Cover = value.Get("cover_image_url").String()
		sourceComic.SourceUri = "web/topic/" + value.Get("id").String()
		sourceComic.Title = value.Get("title").String()
		var Category []string
		for _, v := range value.Get("category").Array() {
			Category = append(Category, v.Str)
		}
		sourceComic.Category = Category
		sourceComic.Author = value.Get("author_name").String()
		sourceComic.LikeCount = value.Get("likes_count").String()
		sourceComic.Popularity = value.Get("popularity").String()
		sourceComic.IsFree = 0
		sourceComic.SourceData = value.String()
		if value.Get("is_free").Bool() == false {
			sourceComic.IsFree = 1
		}

		err := orm.Eloquent.Where("source = ? and source_id = ?", 1, id).FirstOrCreate(&sourceComic).Error
		if err != nil {
			panic(err)
		}
		rd.RPush(model.SourceComicTASK, sourceComic.Id)
		return true
	})
}
