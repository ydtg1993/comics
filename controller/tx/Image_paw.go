package tx

import (
	"comics/common"
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/tebeka/selenium"
	"math"
	"strconv"
	"time"
)

func ImagePaw() {
	rob := robot.GetRob()
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)
	rob.WebDriver.ResizeWindow("", 1400, 900)
	taskLimit := 50
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}
		var sourceChapter model.SourceChapter
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceChapter).Error; err != nil {
			logs.Info("未找到chapter id=" + id)
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		imgContain, err := rob.WebDriver.FindElement(selenium.ByClassName, "comic-contain")
		if err != nil {
			logs.Error(fmt.Sprintf("未找到图片列表Dom: comic-contain source = %d chapter_url = %s err = %s",
				config.Spe.SourceId,
				sourceChapter.SourceUrl, err.Error()))
			robot.ReSetUp(config.Spe.Maxthreads)
			return
		}
		var arg []interface{}
		wait := 10
		vh, err := rob.WebDriver.ExecuteScript(`return document.getElementById("comicContain").clientHeight`, arg)
		if err == nil {
			vhi, err := strconv.Atoi(tools.UnknowToString(vh))
			if err == nil {
				wait = int(math.Ceil(float64(vhi) / float64(1000)))
			}
		}
		rob.WebDriver.ExecuteScript(`
let f = setInterval(toBottom,500);
function toBottom(){
 let dom = document.getElementById("mainView")
 const currentScroll = dom.scrollTop 
 const clientHeight = dom.clientHeight; 
 const scrollHeight = dom.scrollHeight; 
 if (scrollHeight - 10 > currentScroll + clientHeight) {
 	 dom.scrollTo({'left':0,'top': currentScroll + 800,behavior: 'smooth'})
  }else{
	 clearInterval(f)
  }
}
`, arg)
		t := time.NewTicker(time.Second * time.Duration(wait))
		<-t.C

		var sourceImage model.SourceImage
		imageList, _ := imgContain.FindElements(selenium.ByTagName, "li")
		for _, image := range imageList {
			class, err := image.GetAttribute("class")
			if err == nil && class == "main_ad_top" {
				continue
			}
			img, err := image.FindElement(selenium.ByTagName, "img")
			if err != nil {
				continue
			}
			source, _ := img.GetAttribute("src")
			sourceImage.SourceData = append(sourceImage.SourceData, source)
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
			file := common.DownFile(img, dir, fmt.Sprintf("%d.jpg", key), ck)
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
