package controller

import (
	"comics/global/orm"
	"comics/model"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/rd"
	"fmt"
	"github.com/tebeka/selenium"
	"sync"
	"time"
)

func ChapterPaw(tunnel int) {
	Rob := robot.Robot{Port: 19993 + tunnel}
	Rob.Start(config.Spe.SourceUrl)

	wg := sync.WaitGroup{}
	taskLimit := 100
	wg.Add(taskLimit)
	for limit := 0; limit < taskLimit; limit++ {
		id, err := rd.LPop(model.SourceComicTASK)
		if err != nil || id == "" {
			for i := 0; i < (taskLimit - limit); i++ {
				wg.Done()
			}
			return
		}

		var sourceComic model.SourceComic
		if err := orm.Eloquent.Where("id = ?", id).First(&sourceComic).Error; err != nil {
			wg.Done()
			continue
		}
		Rob.WebDriver.Get(config.Spe.SourceUrl + "/" + sourceComic.SourceUri)
		t := time.NewTicker(time.Second * 1)
		<-t.C
		var arg []interface{}
		Rob.WebDriver.ExecuteScript("window.scrollBy(0,10000)", arg)
		listElements, err := Rob.WebDriver.FindElements(selenium.ByClassName, "TopicItem")
		if err != nil {
			continue
		}
		var itemElement selenium.WebElement
		for _, itemElement = range listElements {
			dom, err := itemElement.FindElement(selenium.ByClassName, "img")
			var (
				title = ""
				img   = ""
				date  = ""
				free  = 0
			)
			if err == nil {
				title, _ = dom.GetAttribute("alt")
				img, _ = dom.GetAttribute("src")
			}
			dom, err = itemElement.FindElement(selenium.ByClassName, "date")
			if err == nil {
				date, _ = dom.Text()
			}
			_, err = itemElement.FindElement(selenium.ByClassName, "lockedIcon")
			if err != nil { //收费
				free = 1
			}
			fmt.Println(title, img, date)
		}
	}
	wg.Wait()
	Rob.Service.Stop()
	Rob.WebDriver.Close()
}
