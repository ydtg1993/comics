package main

import (
	"comics/controller"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/database"
	"comics/tools/log"
	"comics/tools/rd"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"os"
	"runtime"
	"strconv"
)

var Source *controller.SourceStrategy

func main() {
	Setup()
	fmt.Println(config.Spe.SourceUrl)
	go controller.TaskComic(Source)

	go controller.TaskChapter(Source)
	go controller.TaskChapterUpdate()

	go controller.TaskImage(Source)

	controller.TaskDownImage()
}

func Setup() {
	err := config.Spe.SetUp()
	if err != nil {
		panic(err)
	}

	url := os.Getenv("SOURCE_URL")
	if url != "" {
		config.Spe.SourceUrl = url
	}
	Source = controller.SourceOperate(config.Spe.SourceUrl)
	config.Spe.RedisDb = config.Spe.SourceId

	mylog := new(log.LogsManage)
	err = mylog.SetUp()
	if err != nil {
		panic(err)
	}

	db := new(database.MysqlManage)
	err = db.Setup()
	if err != nil {
		panic(err)
	}

	redisManage := new(rd.RedisManage)
	err = redisManage.SetUp()
	if err != nil {
		panic(err)
	}
	controller.TaskStepRecord += strconv.Itoa(config.Spe.SourceId)
	rd.Delete(controller.TaskStepRecord)

	if config.Spe.SeleniumPath != "" {
		go robot.SetUp()
	}
	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())
}
