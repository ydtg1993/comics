package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tebeka/selenium"
	"strconv"
	"time"
)

func ImagePaw() {
	rob := robot.GetRob()
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)

	taskLimit := 100
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}
		var sourceChapter model.SourceChapter
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceChapter).Error; err != nil {
			logs.Info(fmt.Sprintf("未找到chapter_id = %s source = %d", id, config.Spe.SourceId))
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		var arg []interface{}
		rob.WebDriver.ExecuteScript("window.scrollBy(0,1000000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C

		imgList, err := rob.WebDriver.FindElements(selenium.ByClassName, "img-box")
		if err != nil {
			logs.Error(fmt.Sprintf("未找到图片列表Dom: imgList source = %d chapter_id = %s err = %s",
				config.Spe.SourceId,
				id, err.Error()))
			robot.ReSetUp(config.Spe.Maxthreads)
			return
		}
		var sourceImage model.SourceImage
		for _, img := range imgList {
			dom, err := img.FindElement(selenium.ByClassName, "img")
			if err != nil {
				logs.Error("无法抓取图片页Dom: img " + sourceChapter.SourceUrl)
			}
			img, err := dom.GetAttribute("data-src")
			if err != nil {
				logs.Error("无法抓取图片页Dom: img attr " + sourceChapter.SourceUrl)
			}
			sourceImage.SourceData = append(sourceImage.SourceData, img)
		}
		sourceImage.ChapterId, _ = strconv.Atoi(id)
		sourceImage.Images = model.Images{}
		cookies, _ := rob.WebDriver.GetCookies()
		download(
			sourceChapter.ComicId,
			&sourceImage,
			cookies, sourceImage.SourceData)

		var exists bool
		orm.Eloquent.Model(model.SourceImage{}).Where("chapter_id = ?", id).First(&exists)
		if exists == false {
			err = orm.Eloquent.Create(&sourceImage).Error
			if err != nil {
				logs.Error(fmt.Sprintf("image数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error()))
			}
		} else {
			err = orm.Eloquent.Model(model.SourceImage{}).Where("chapter_id = ?", id).Updates(map[string]interface{}{
				"images":      sourceImage.Images,
				"source_data": sourceImage.SourceData,
				"state":       sourceImage.State,
			}).Error
			if err != nil {
				logs.Error(fmt.Sprintf("image数据更新失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error()))
			}
		}
	}
}

func download(comicId int, sourceImage *model.SourceImage, cookies []selenium.Cookie, images model.Images) {
	ck := make(map[string]string)
	for _, cookie := range cookies {
		ck[cookie.Name] = cookie.Value
	}
	for key, img := range images {
		dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d", comicId, sourceImage.ChapterId)

		state := 0
		for i := 0; i < 3; i++ {
			file := common.DownFile(img, dir, fmt.Sprintf("%d.webp", key), ck)
			if file != "" {
				state = 1
				sourceImage.Images = append(sourceImage.Images, file)
				break
			}
		}
		if state == 0 {
			sourceImage.Images = model.Images{}
			sourceImage.State = state
			return
		}
		sourceImage.State = state
	}
}
