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
	"github.com/tebeka/selenium"
	"math"
	"regexp"
	"strconv"
	"time"
)

func ImagePaw() {
	rob := robot.GetRob([]int{})
	if rob == nil {
		return
	}
	robot.ResetRob(rob)
	defer func() {
		rob.State = 0
		rob.Lock.Unlock()
	}()

	taskLimit := 10
	for limit := 0; limit < taskLimit; limit++ {
		common.StopSignal("图片任务挂起")
		id, err := rd.LPop(common.SourceChapterTASK)
		if err != nil || id == "" {
			return
		}
		sourceChapter := new(model.SourceChapter)
		if orm.Eloquent.Where("id = ?", id).First(&sourceChapter); sourceChapter.Id == 0 {
			continue
		}
		rob.WebDriver.Get(sourceChapter.SourceUrl)
		sourceImage := new(model.SourceImage)
		sourceImage.Source = config.Spe.SourceId
		sourceImage.ComicId = sourceChapter.ComicId
		sourceImage.ChapterId = sourceChapter.Id
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}
		browserList(rob, sourceImage, sourceChapter)
		if len(sourceImage.SourceData) == 0 {
			msg := fmt.Sprintf("未找到图片资源列表: source = %d comic_id = %d chapter_url = %s",
				config.Spe.SourceId,
				sourceChapter.ComicId,
				sourceChapter.SourceUrl)
			rob.WebDriver.Refresh()
			model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
			rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			continue
		}
		record := new(model.SourceImage)
		exists := record.Exists(sourceChapter.Id)
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
			record.Images = sourceImage.Images
			record.SourceData = sourceImage.SourceData
			record.State = sourceImage.State
			err = orm.Eloquent.Save(record).Error
			if err != nil {
				msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceChapterId,
					err.Error())
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片数据更新错误", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, sourceImage.Id)
			}
		}
	}
}

func browserList(rob *robot.Robot, sourceImage *model.SourceImage, sourceChapter *model.SourceChapter) {
	script := `
if (document.getElementById("mainView").scrollTop == 0){
		let f1 = setInterval(()=>{
	 let dom = document.getElementById("mainView")
	 const currentScroll = dom.scrollTop 
	 const clientHeight = dom.clientHeight; 
	 const scrollHeight = dom.scrollHeight; 
	 if (scrollHeight - 10 > currentScroll + clientHeight) {
		 dom.scrollTo({'left':0,'top': currentScroll + 1200,behavior: 'smooth'})
	  }else{
		 clearInterval(f1)
	  }
	},500);
}else{
	let f2 = setInterval(()=>{
	 let dom = document.getElementById("mainView")
	 const currentScroll = dom.scrollTop 
	 const clientHeight = dom.clientHeight; 
	 const scrollHeight = dom.scrollHeight; 
	 if (scrollHeight + 50 > currentScroll + clientHeight) {
		 dom.scrollTo({'left':0,'top': currentScroll - 1400,behavior: 'smooth'})
	  }else{
		 clearInterval(f2)
	  }
	},500);
}`

	for tryLimit := 0; tryLimit <= 6; tryLimit++ {
		imgContain, err := rob.WebDriver.FindElement(selenium.ByClassName, "comic-contain")
		if err != nil {
			if tryLimit > 3 {
				if tryLimit == 6 {
					msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %s",
						config.Spe.SourceId,
						sourceChapter.ComicId,
						sourceChapter.SourceUrl)
					model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
					rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
					return
				}
				robot.ResetRob(rob)
				rob.WebDriver.Get(sourceChapter.SourceUrl)
			}
			continue
		}
		var arg []interface{}
		wait := 30
		vh, err := rob.WebDriver.ExecuteScript(`return document.getElementById("comicContain").clientHeight`, arg)
		if err == nil {
			vhi, err := strconv.Atoi(tools.UnknowToString(vh))
			if err == nil {
				wait = int(math.Ceil(float64(vhi) / float64(2000)))
			}
		}
		rob.WebDriver.ExecuteScript(script, arg)
		t := time.NewTicker(time.Second * time.Duration(wait))
		<-t.C

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
			match, _ := regexp.MatchString("pixel.gif", source)
			if source != "" && match != true {
				sourceImage.SourceData = append(sourceImage.SourceData, source)
			} else {
				sourceImage.SourceData = model.Images{}
				if tryLimit%2 == 0 {
					rob.WebDriver.Refresh()
				}
				break
			}
		}
		if len(sourceImage.SourceData) > 0 {
			return
		}
	}
}
