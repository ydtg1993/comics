package tx

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
			logs.Info("未找到chapter id=" + id)
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		var arg []interface{}
		rob.WebDriver.ExecuteScript("window.scrollBy(0,1000000)", arg)
		t := time.NewTicker(time.Second * 2)
		<-t.C

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
