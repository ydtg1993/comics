package kk

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
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
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		var arg []interface{}
		rob.WebDriver.ExecuteScript("window.scrollBy(0,1000000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C

		imgList, err := rob.WebDriver.FindElements(selenium.ByClassName, "img-box")
		if err != nil {
			msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %d chapter_url = %s err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.Id,
				sourceChapter.SourceUrl,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片列表未找到 重启机器人", 3)
			return
		}
		var sourceImage model.SourceImage
		sourceImage.Images = model.Images{}
		for _, img := range imgList {
			dom, err := img.FindElement(selenium.ByClassName, "img")
			if err == nil {
				img, err := dom.GetAttribute("data-src")
				if err == nil {
					sourceImage.SourceData = append(sourceImage.SourceData, img)
				}
			}
		}
		if len(sourceImage.SourceData) == 0 {
			msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %d chapter_url = %s err = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.Id,
				sourceChapter.SourceUrl,
				err.Error())
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片列表未找到", 3)
			continue
		}
		sourceImage.ChapterId, _ = strconv.Atoi(id)
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
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误", 3)
			}
		} else {
			err = orm.Eloquent.Model(model.SourceImage{}).Where("chapter_id = ?", id).Updates(map[string]interface{}{
				"images":      sourceImage.Images,
				"source_data": sourceImage.SourceData,
				"state":       sourceImage.State,
			}).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据更新失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片数据更新错误", 3)
			}
		}
	}
}

func download(comicId int, sourceImage *model.SourceImage, cookies []selenium.Cookie, images model.Images) {
	ck := make(map[string]string)
	for _, cookie := range cookies {
		ck[cookie.Name] = cookie.Value
	}
	dir := fmt.Sprintf(config.Spe.DownloadPath+"chapter/%d/%d", comicId, sourceImage.ChapterId)
	for key, img := range images {
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
