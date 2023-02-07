package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/tebeka/selenium"
	"path/filepath"
	"strconv"
	"time"
)

func ChapterPaw() {
	rob := robot.GetRob([]int{0, 2, 4})
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)

	taskLimit := 50
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(common.SourceComicTASK)
		if err != nil || id == "" {
			return
		}

		var sourceComic model.SourceComic
		if orm.Eloquent.Where("id = ?", id).First(&sourceComic); sourceComic.Id == 0 {
			continue
		}
		rob.WebDriver.Get(sourceComic.SourceUrl)
		t := time.NewTicker(time.Second * 2)
		<-t.C

		_, check := rob.WebDriver.FindElement(selenium.ByXPATH, "//*[@id='special_bg']")
		if check != nil {
			msg := fmt.Sprintf("未找到漫画详情页内容 source = %d comic_id = %s comic_url = %s",
				config.Spe.SourceId, id, sourceComic.SourceUrl)
			model.RecordFail(sourceComic.SourceUrl, msg, "漫画详情未找到", 2)
			rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			continue
		}
		var arg []interface{}
		rob.WebDriver.ExecuteScript(`
let mouseoverEvent = new Event('mouseover');
document.getElementsByClassName("chapter-page-btn-all")[0].dispatchEvent(mouseoverEvent);
`, arg)
		elms, _ := rob.WebDriver.FindElements(selenium.ByCSSSelector, ".chapter-page-more>a")
		if len(elms) == 0 {
			getChapter(&sourceComic, rob)
		}
		for _, elem := range elms {
			elem.Click()
			t := time.NewTicker(time.Second * 1)
			<-t.C
			getChapter(&sourceComic, rob)
		}

		detail, err := rob.WebDriver.FindElement(selenium.ByCSSSelector, ".works-intro-short")
		if err == nil {
			sourceComic.Description, _ = detail.Text()
		}
		tags, err := rob.WebDriver.FindElements(selenium.ByCSSSelector, ".works-intro-tags-item")
		if err == nil && len(tags) > 0 {
			for _, tag := range tags {
				tagString, err := tag.Text()
				if err == nil {
					sourceComic.Label = append(sourceComic.Label, tagString)
				}
			}
		}
		like, err := rob.WebDriver.FindElement(selenium.ByCSSSelector, "#redcount")
		if err == nil {
			sourceComic.LikeCount, _ = like.Text()
		}
		sourceComic.Region = "国漫"
		var total int64
		orm.Eloquent.Model(model.SourceChapter{}).Where("comic_id = ?", sourceComic.Id).Count(&total)
		sourceComic.ChapterCount = int(total)
		orm.Eloquent.Save(&sourceComic)
	}
}

func getChapter(sourceComic *model.SourceComic, rob *robot.Robot) {
	elms, err := rob.WebDriver.FindElements(selenium.ByCSSSelector, ".chapter-page-all .works-chapter-item")
	if err != nil {
		return
	}
	if len(elms) == 0 {
		return
	}
	for sort, elem := range elms {
		sourceChapter := new(model.SourceChapter)
		sourceChapter.ComicId = sourceComic.Id
		sourceChapter.Source = config.Spe.SourceId
		sourceChapter.Sort = sort
		_, pay := elem.FindElement(selenium.ByTagName, "ui-icon-pay")
		if pay == nil {
			sourceChapter.IsFree = 1
		}
		a, err := elem.FindElement(selenium.ByTagName, "a")
		if err != nil {
			continue
		}
		sourceChapter.Title, _ = a.GetAttribute("title")
		url, err := a.GetAttribute("href")
		if err == nil {
			sourceChapter.SourceChapterId, _ = strconv.Atoi(filepath.Base(url))
			sourceChapter.SourceUrl = url
		}
		if sourceChapter.SourceChapterId == 0 {
			msg := fmt.Sprintf("章节id没有查找到 source = %d comic_id = %d chapter_url = %s",
				config.Spe.SourceId,
				sourceComic.Id, sourceChapter.SourceUrl)
			model.RecordFail(sourceComic.SourceUrl, msg, "漫画章节未找到", 2)
			rd.RPush(common.SourceComicRetryTask, sourceComic.Id)
			continue
		}

		var exists bool
		orm.Eloquent.Model(model.SourceChapter{}).Select("count(*) > 0").Where("source = ? and source_url = ?",
			config.Spe.SourceId,
			sourceChapter.SourceUrl).Find(&exists)
		if exists == false {
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
			}
		}
	}
}
