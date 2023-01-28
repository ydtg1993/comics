package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/tools/config"
	"comics/tools/rd"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"strings"
)

func ChapterPaw() {
	for limit := 0; limit < 100; limit++ {
		id, err := rd.LPop(model.SourceComicTASK)
		if err != nil {
			logs.Error("source:comic:task 进程", "redis 读取错误:", err.Error())
			return
		}
		if id == "" {
			return
		}

		var sourceComic model.SourceComic
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceComic).Error; err != nil {
			continue
		}

		content, err := requestHtml(config.Spe.SourceUrl+"/"+sourceComic.SourceUri, 3)
		if err != nil {
			return
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			logs.Error("无法抓取目漫画详情页Dom: " + sourceComic.SourceUri)
			return
		}
		sourceComicExt := new(model.SourceComicExt)
		sourceComicExt.ComicId, _ = strconv.Atoi(id)
		sourceComicExt.Author = doc.Find("div.nickname").Text()

	}
}
