package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"github.com/tebeka/selenium"
	"time"
)

func ChapterPaw() {
	var Rob *robot.Robot
	for _, robot := range robot.Swarm {
		if robot.State == 1 {
			continue
		}
		robot.Lock.Lock()
		if robot.State == 1 {
			continue
		}
		robot.State = 1
		Rob = robot
		break
	}
	if Rob == nil {
		return
	}
	defer func() {
		Rob.State = 0
		Rob.Lock.Unlock()
	}()

	taskLimit := 10
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceComicTASK)
		if err != nil || id == "" {
			return
		}

		var sourceComic model.SourceComic
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceComic).Error; err != nil {
			continue
		}
		Rob.WebDriver.Get("https://" + config.Spe.SourceUrl + "/" + sourceComic.SourceUri)
		var arg []interface{}
		Rob.WebDriver.ExecuteScript("window.scrollBy(0,10000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C
		listElements, err := Rob.WebDriver.FindElements(selenium.ByClassName, "TopicItem")
		if err != nil {
			continue
		}

		for sort, itemElement := range listElements {
			dom, err := itemElement.FindElement(selenium.ByClassName, "img")
			sourceChapter := new(model.SourceChapter)
			sourceChapter.Source = 1
			sourceChapter.ComicId = sourceComic.SourceId
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
					sourceChapter.SourceUri, err = dom.GetAttribute("href")
					if err == nil {
						sourceChapter.SourceChapterId = tools.FindStringNumber(sourceChapter.SourceUri)
					}
				}
			}
			if sourceChapter.SourceChapterId > 0 {
				orm.Eloquent.Where("source = ? and source_id = ? and source_chapter_id = ?",
					1,
					sourceComic.SourceId,
					sourceChapter.SourceChapterId).FirstOrCreate(&sourceChapter)
				rd.RPush(model.SourceChapterTASK, sourceChapter.Id)
			}
		}
	}
}
