package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
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
	url := fmt.Sprintf(config.Spe.SourceUrl+"/tag/%d?region=%d&pays=%d&state=%d&sort=1&page=1",
		tagId, regionId, payId, stateId)
	content, err := requestHtml(url, 3)
	if err != nil {
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		logs.Error("无法抓取分类列表页Dom:" + url)
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
	url := fmt.Sprintf(config.Spe.SourceUrl+"/search/mini/topic/multi_filter?tag_id=%d&label_dimension_origin=%d&pay_status=%d&update_status=%d&sort=1&page=%d&size=48",
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

		orm.Eloquent.Where("source = ? and source_id = ?", 1, id).FirstOrCreate(&sourceComic)

		sourceComicExt := new(model.SourceComicExt)
		sourceComicExt.ComicId = id
		var Category []string
		for _, v := range value.Get("category").Array() {
			Category = append(Category, v.Str)
		}
		sourceComicExt.Category = Category
		sourceComicExt.Author = value.Get("author_name").String()
		sourceComicExt.LikeCount, _ = strconv.Atoi(value.Get("likes_count").String())
		sourceComicExt.Popularity, _ = strconv.Atoi(value.Get("popularity").String())
		sourceComicExt.IsFree = 0
		sourceComicExt.SourceData = content.String()
		if value.Get("is_free").Bool() == false {
			sourceComicExt.IsFree = 1
		}
		err := orm.Eloquent.Where("comic_id = ?", id).FirstOrCreate(&sourceComicExt).Error
		if err != nil {
			panic(err)
		}
		rd.RPush(model.SourceComicTASK, sourceComic.Id)
		return true
	})
}
