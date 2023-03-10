package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/tebeka/selenium"
	"time"
)

func ChapterPaw() {
	rob := robot.GetRob([]int{0, 1})
	if rob == nil {
		return
	}
	robot.ResetRob(rob)
	defer func() {
		rob.State = 0
		rob.Lock.Unlock()
	}()

	taskLimit := 30
	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("章节任务挂起")
		id, err := rd.LPop(common.SourceComicTASK)
		if err != nil || id == "" {
			return
		}

		var sourceComic model.SourceComic
		if orm.Eloquent.Where("id = ?", id).First(&sourceComic); sourceComic.Id == 0 {
			continue
		}
		for tryLimit := 0; tryLimit <= 5; tryLimit++ {
			rob.WebDriver.Get(sourceComic.SourceUrl)
			var arg []interface{}
			rob.WebDriver.ExecuteScript("window.scrollTo({'left':0,'top': 10000000,behavior: 'smooth'})", arg)
			t := time.NewTicker(time.Second * 3)
			<-t.C
			listElements, err := rob.WebDriver.FindElements(selenium.ByClassName, "TopicItem")
			if err != nil {
				if tryLimit > 3 {
					if tryLimit == 5 {
						msg := fmt.Sprintf("章节列表未找到 source = %d comic_id = %s comic_url = %s err = %s",
							config.Spe.SourceId,
							id, sourceComic.SourceUrl, err.Error())
						model.RecordFail(sourceComic.SourceUrl, msg, "章节列表未找到", 2)
						rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
					}
					robot.ResetRob(rob)
				}
				continue
			}
			chapterList(&sourceComic, listElements)
			break
		}

		detail, err := rob.WebDriver.FindElement(selenium.ByClassName, "detailsBox")
		if err == nil {
			sourceComic.Description, _ = detail.Text()
		}
		var total int64
		orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
		sourceComic.ChapterCount = int(total)
		orm.Eloquent.Save(&sourceComic)
	}
}

func chapterList(sourceComic *model.SourceComic, listElements []selenium.WebElement) {
	for sort, itemElement := range listElements {
		if sourceComic.Retry == 0 && sort < sourceComic.ChapterPick {
			continue
		}
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
					if sourceChapter.SourceUrl == "javascript:void(0);" {
						sourceChapter.SourceUrl = ""
						sourceChapter.SourceChapterId = 0
					} else {
						sourceChapter.SourceChapterId = tools.FindStringNumber(sourceChapter.SourceUrl)
					}
				}
			}
		}

		if sourceChapter.SourceChapterId == 0 {
			if sourceChapter.IsFree == 1 {
				msg := fmt.Sprintf("章节还没有完成购买 source = %d comic_id = %d chapter_url = %s chapter_name = %s",
					config.Spe.SourceId,
					sourceComic.Id, sourceChapter.SourceUrl,
					sourceChapter.Title)
				model.RecordFail(sourceChapter.SourceUrl, msg, "章节没有购买", 2)
				//rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			} else {
				msg := fmt.Sprintf("章节id没有查找到 source = %d comic_id = %d chapter_url = %s chapter_name = %s",
					config.Spe.SourceId,
					sourceComic.Id, sourceChapter.SourceUrl,
					sourceChapter.Title)
				model.RecordFail(sourceChapter.SourceUrl, msg, "章节id没有查找到", 2)
				rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			}
			continue
		}

		exists := new(model.SourceChapter).Exists(sourceComic.Id, sourceChapter.SourceUrl)
		if exists == false {
			sourceComic.ChapterPick = sort
			err = orm.Eloquent.Create(&sourceChapter).Error
			if err != nil {
				msg := fmt.Sprintf("chapter数据导入失败 source = %d comic_id = %d chapter_url = %s err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl,
					err.Error())
				model.RecordFail(sourceComic.SourceUrl, msg, "漫画章节入库错误", 2)
				rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			} else {
				rd.RPush(common.SourceChapterTASK, sourceChapter.Id)
				sourceComic.LastChapterUpdateAt = time.Now()
			}
		}
	}
}
