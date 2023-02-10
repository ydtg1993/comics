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
	"strconv"
	"time"
)

func ImagePaw() {
	rob := robot.GetRob([]int{})
	if rob == nil {
		return
	}
	defer robot.ResetRob(rob)

	taskLimit := 50
	for limit := 0; limit < taskLimit; limit++ {
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
		sourceImage.ChapterId = sourceChapter.Id
		sourceImage.Images = model.Images{}
		sourceImage.SourceData = model.Images{}
		for tryLimit := 0; tryLimit <= 3; tryLimit++ {
			imgContain, err := rob.WebDriver.FindElement(selenium.ByClassName, "comic-contain")
			if err != nil && tryLimit == 3 {
				msg := fmt.Sprintf("未找到图片列表: source = %d comic_id = %d chapter_url = %s",
					config.Spe.SourceId,
					sourceChapter.ComicId,
					sourceChapter.SourceUrl)
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片资源未找到", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
				rob.WebDriver.Refresh()
				continue
			}
			var arg []interface{}
			wait := 30
			vh, err := rob.WebDriver.ExecuteScript(`return document.getElementById("comicContain").clientHeight`, arg)
			if err == nil {
				vhi, err := strconv.Atoi(tools.UnknowToString(vh))
				if err == nil {
					wait = int(math.Ceil(float64(vhi) / float64(3000)))
				}
			}
			rob.WebDriver.ExecuteScript(`
if (document.getElementById("mainView").scrollTop == 0){
		let f1 = setInterval(()=>{
	 let dom = document.getElementById("mainView")
	 const currentScroll = dom.scrollTop 
	 const clientHeight = dom.clientHeight; 
	 const scrollHeight = dom.scrollHeight; 
	 if (scrollHeight - 10 > currentScroll + clientHeight) {
		 dom.scrollTo({'left':0,'top': currentScroll + 1600,behavior: 'smooth'})
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
		 dom.scrollTo({'left':0,'top': currentScroll - 1600,behavior: 'smooth'})
	  }else{
		 clearInterval(f2)
	  }
	},500);
}
`, arg)
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
				if source != "" {
					sourceImage.SourceData = append(sourceImage.SourceData, source)
				}
			}
			if len(sourceImage.SourceData) > 0 {
				break
			}
		}

		var exists bool
		msg := fmt.Sprintf("图片数据导入失败 source = %d comic_id = %d chapter_id = %d err = %s",
			config.Spe.SourceId,
			sourceChapter.ComicId,
			sourceChapter.SourceChapterId,
			err.Error())
		orm.Eloquent.Model(model.SourceImage{}).Select("id > 0").Where("chapter_id = ?", id).First(&exists)
		if exists == false {
			err = orm.Eloquent.Create(&sourceImage).Error
			if err != nil {
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
				model.RecordFail(sourceChapter.SourceUrl, msg, "图片数据更新错误", 3)
				rd.RPush(common.SourceChapterRetryTask, sourceChapter.Id)
			} else {
				rd.RPush(common.SourceImageTASK, sourceImage.Id)
			}
		}
	}
}
