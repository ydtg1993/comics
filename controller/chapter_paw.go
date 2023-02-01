package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tebeka/selenium"
	"time"
)

func ChapterPaw() {
	rob := robot.GetRob()
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)

	taskLimit := 100
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceComicTASK)
		if err != nil || id == "" {
			return
		}

		var sourceComic model.SourceComic
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceComic).Error; err != nil {
			logs.Info("未找到comic id=" + id)
			continue
		}
		rob.WebDriver.Get(sourceComic.SourceUrl)
		var arg []interface{}
		rob.WebDriver.ExecuteScript("window.scrollBy(0,100000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C
		listElements, err := rob.WebDriver.FindElements(selenium.ByClassName, "TopicItem")
		if err != nil {
			continue
		}

		for sort, itemElement := range listElements {
			dom, err := itemElement.FindElement(selenium.ByClassName, "img")
			sourceChapter := new(model.SourceChapter)
			sourceChapter.Source = 1
			sourceChapter.ComicId = sourceComic.Id
			sourceChapter.Sort = sort
			if err == nil {
				sourceChapter.Title, err = dom.GetAttribute("alt")
				sourceChapter.Cover, _ = dom.GetAttribute("src")
			} else {
				dom, err = itemElement.FindElement(selenium.ByClassName, "imgCover")
				if err == nil {
					sourceChapter.Title, err = dom.GetAttribute("alt")
					sourceChapter.Cover, _ = dom.GetAttribute("src")
				}
			}

			dom, err = itemElement.FindElement(selenium.ByClassName, "title")
			if err == nil {
				_, err = dom.FindElement(selenium.ByClassName, "lockedIcon")
				if err == nil { //收费
					sourceChapter.IsFree = 1
				}
				dom, err = dom.FindElement(selenium.ByTagName, "a")
				if err == nil {
					sourceChapter.SourceUrl, err = dom.GetAttribute("href")
					if err == nil {
						sourceChapter.SourceChapterId = tools.FindStringNumber(sourceChapter.SourceUrl)
					}
				}
			}
			if sourceChapter.SourceChapterId > 0 {
				var exists bool
				orm.Eloquent.Model(model.SourceChapter{}).Select("count(*) > 0").Where("source = ? and comic_id = ? and source_chapter_id = ?",
					1,
					sourceComic.Id,
					sourceChapter.SourceChapterId).Find(&exists)
				if exists == false {
					err = orm.Eloquent.Create(&sourceChapter).Error
					if err != nil {
						logs.Error(fmt.Sprintf("chapter数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
							1,
							sourceChapter.ComicId,
							sourceChapter.SourceChapterId,
							err.Error()))
					} else {
						rd.RPush(model.SourceChapterTASK, sourceChapter.Id)
					}
				}
			}
		}
	}
}
