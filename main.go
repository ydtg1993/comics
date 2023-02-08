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
	"sync"
	"time"
)

var Source *controller.SourceStrategy
var TaskStepRecord = "task:step:record:"

func main() {
	Setup()

	go TaskComic(Source)

	go TaskChapter(Source)

	TaskImage(Source)
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
	TaskStepRecord += strconv.Itoa(config.Spe.SourceId)
	rd.Delete(TaskStepRecord)

	go robot.SetUp()
	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())
}

func TaskComic(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Hour * 36)
	defer t.Stop()
	rd.RPush(TaskStepRecord, fmt.Sprintf("漫画-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
	source.ComicPaw()

	for {
		<-t.C
		rd.Delete(TaskStepRecord)
		rd.RPush(TaskStepRecord, fmt.Sprintf("漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		source.ComicUpdate()
	}
}

func TaskChapter(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Minute * 3)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(config.Spe.Maxthreads)
		rd.RPush(TaskStepRecord, fmt.Sprintf("章节-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				source.ChapterPaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskChapterUpdate(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Hour * 11)
	defer t.Stop()
	for {
		<-t.C
		rd.RPush(TaskStepRecord, fmt.Sprintf("连载漫画更新-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		source.ComicUpdate()
	}
}

func TaskImage(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Minute * 3)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(config.Spe.Maxthreads)
		rd.RPush(TaskStepRecord, fmt.Sprintf("图片-进程开始 %s %s", config.Spe.SourceUrl, time.Now().String()))
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				source.ImagePaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
