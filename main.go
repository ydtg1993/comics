package main

import (
	"comics/controller"
	"comics/robot"
	"comics/tools/config"
	"comics/tools/database"
	"comics/tools/log"
	"comics/tools/rd"
	"github.com/beego/beego/v2/core/logs"
	"runtime"
	"sync"
	"time"
)

func main() {
	Setup()

	go TaskComic()

	go TaskChapter()

	TaskImage()
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

func TaskComic() {
	t := time.NewTicker(time.Hour * 6)
	defer t.Stop()

	controller.ComicPaw()

	for {
		<-t.C
		controller.ComicUpdate()
	}
}

func TaskChapter() {
	t := time.NewTicker(time.Second * 15)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(config.Spe.Maxthreads)
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				controller.ChapterPaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func TaskImage() {
	t := time.NewTicker(time.Second * 20)
	defer t.Stop()
	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(config.Spe.Maxthreads)
		for i := 0; i < config.Spe.Maxthreads; i++ {
			go func() {
				controller.ImagePaw()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
