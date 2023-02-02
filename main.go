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
	"runtime"
	"strconv"
	"sync"
	"time"
)

var TaskStepRecord = "task:step:record:"

func main() {
	Setup()

	TaskStepRecord += strconv.Itoa(config.Spe.SourceId)
	rd.Delete(TaskStepRecord)
	source := controller.SourceOperate(config.Spe.SourceUrl)

	//go TaskComic(source)

	//go TaskChapter(source)

	TaskImage(source)
}

func Setup() {
	err := config.Spe.SetUp()
	if err != nil {
		panic(err)
	}

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

	go robot.SetUp(config.Spe.Maxthreads)

	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())
}

func TaskComic(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Hour * 6)
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
	t := time.NewTicker(time.Second * 5)
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

func TaskImage(source *controller.SourceStrategy) {
	t := time.NewTicker(time.Second * 1)
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
