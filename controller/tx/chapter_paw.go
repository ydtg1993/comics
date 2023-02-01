package tx

import (
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/rd"
	"github.com/beego/beego/v2/core/logs"
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

	}
}
