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
)

func ImagePaw() {
	rob := robot.GetRob([]int{})
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)

	taskLimit := 50
	for limit := 0; limit < taskLimit; limit++ {
		signal := common.Signal("图片")
		if signal == true {
			return
		}
		id, err := rd.LPop(common.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}
		var sourceChapter model.SourceChapter
		if orm.Eloquent.Where("id = ?", id).First(&sourceChapter); sourceChapter.Id == 0 {
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		sourceImage := new(model.SourceImage)
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}
		sourceImage.ChapterId = sourceChapter.Id
		for tryLimit := 0; tryLimit <= 3; tryLimit++ {
			var arg []interface{}
			rob.WebDriver.ExecuteScript("window.scrollBy(0,1000000)", arg)
			imgList, err := rob.WebDriver.FindElements(selenium.ByClassName, "img-box")
			if err != nil {
				msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %d chapter_url = %s err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.Id,
					sourceChapter.SourceUrl,
					err.Error())
				if tryLimit == 2 {
					model.RecordFail(sourceChapter.SourceUrl, msg, "图片列表未找到", 3)
					rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				}
				continue
			}
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
				msg := fmt.Sprintf("未找到图片资源列表: source = %d comic_id = %d chapter_url = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl)
				rob.WebDriver.Refresh()
				if tryLimit == 3 {
					model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
					rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				}
				continue
			} else {
				break
			}
		}

		if len(sourceImage.SourceData) == 0 {
			continue
		}
		var exists bool
		orm.Eloquent.Model(model.SourceImage{}).Select("id > 0").Where("chapter_id = ?", id).First(&exists)
		if exists == false {
			err = orm.Eloquent.Create(&sourceImage).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, sourceImage.Id)
			}
		} else {
			err = orm.Eloquent.Model(model.SourceImage{}).Where("chapter_id = ?", id).Updates(map[string]interface{}{
				"images":      sourceImage.Images,
				"source_data": sourceImage.SourceData,
				"state":       sourceImage.State,
			}).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片入库错误.更新", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, sourceImage.Id)
			}
		}
	}
}
