package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/gocolly/colly"
	"path/filepath"
	"strconv"
	"time"
)

func ChapterPaw() {
	taskLimit := 50

	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("章节任务挂起")
		id, err := rd.LPop(common.SourceComicTASK)
		if err != nil || id == "" {
			return
		}
		bot := robot.GetColly()
		sourceComic := new(model.SourceComic)
		if orm.Eloquent.Where("id = ?", id).First(&sourceComic); sourceComic.Id == 0 {
			continue
		}

		bot.OnHTML("ol.chapter-page-all", func(e *colly.HTMLElement) {
			e.ForEach(".works-chapter-item", func(sort int, e *colly.HTMLElement) {
				if sourceComic.Retry == 0 && sort < sourceComic.ChapterPick {
					return
				}
				dom := e.DOM.Find("a")
				title, _ := dom.Attr("title")
				url, _ := dom.Attr("href")

				sourceChapter := new(model.SourceChapter)
				sourceChapter.ComicId = sourceComic.Id
				sourceChapter.Source = config.Spe.SourceId
				sourceChapter.Sort = sort
				sourceChapter.Title = title
				sourceChapter.SourceUrl = "https://" + config.Spe.SourceUrl + url
				sourceChapter.SourceChapterId, _ = strconv.Atoi(filepath.Base(url))
				if url == "" || sourceChapter.SourceChapterId == 0 {
					return
				}
				pay := e.DOM.Find("i.ui-icon-pay").Index()
				if pay != -1 {
					sourceChapter.IsFree = 1
				}
				app := e.DOM.Find("span.in-app").Index()
				if app != -1 {
					return
				}
				exists := new(model.SourceChapter).Exists(sourceComic.Id, sourceChapter.SourceUrl)
				if exists == false {
					sourceComic.ChapterPick = sort
					err := orm.Eloquent.Create(&sourceChapter).Error
					if err != nil {
						msg := fmt.Sprintf("chapter数据导入失败 source = %d comic_id = %d chapter_url = %s err = %s",
							config.Spe.SourceId, sourceChapter.ComicId, sourceChapter.SourceUrl, err.Error())
						model.RecordFail(sourceComic.SourceUrl, msg, "漫画章节入库错误", 2)
						rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
					} else {
						rd.RPush(common.SourceChapterTASK, sourceChapter.Id)
						sourceComic.LastChapterUpdateAt = time.Now()
					}
				}
			})
		})

		bot.OnHTML("div.works-intro-wr", func(e *colly.HTMLElement) {
			state := e.DOM.Find("label.works-intro-status").Text()
			description := e.DOM.Find(".works-intro-short").Text()
			like := e.DOM.Find("#coll_count").Text()
			if state == "已完结" {
				sourceComic.IsFinish = 1
			}
			sourceComic.Description = description
			sourceComic.Like = like
			sourceComic.Region = "国漫"
			sourceComic.SourceData, _ = e.DOM.Html()
			var total int64
			orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
			sourceComic.ChapterCount = int(total)
			orm.Eloquent.Save(&sourceComic)
		})

		for i := 0; i <= 3; i++ {
			err := bot.Visit(sourceComic.SourceUrl)
			if err != nil {
				bot = robot.GetColly()
				if i == 3 {
					model.RecordFail(sourceComic.SourceUrl, "无法获取漫画详情 :"+sourceComic.SourceUrl, "漫画详情错误", 2)
					rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
				}
			} else {
				break
			}
		}
	}
}
