package controller

import (
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

	taskLimit := 50
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}
		var sourceChapter model.SourceChapter
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceChapter).Error; err != nil {
			fmt.Println(err.Error())
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		var arg []interface{}
		rob.WebDriver.ExecuteScript("window.scrollBy(0,1000000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C

		imgList, err := rob.WebDriver.FindElements(selenium.ByClassName, "img-box")
		if err != nil {
			logs.Error("无法抓取图片页Dom: imgList " + sourceChapter.SourceUrl)
			continue
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

		orm.Eloquent.Where("chapter_id = ?", id).FirstOrCreate(&sourceImage)

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
			file := DownFile(img, dir, fmt.Sprintf("%d.webp", key), ck)
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
