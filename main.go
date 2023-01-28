package main

import (
	"comics/controller"
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

	TaskComic()

	//TaskChapter()

	//go TaskImage()

	/*done := make(chan int)
	Robot := new(robot.Robot)
	defer Robot.Service.Stop()
	Robot.Start("https://ac.qq.com/Comic/all/page/1")
	Robot.TapIn(selenium.ByXPATH, "/html/body/div[3]/div[2]/div/div[2]/ul/li[1]/div[1]/a")
	<-done*/
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

	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())
}

func TaskComic() {
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()

	controller.ComicPaw()

	for {
		<-t.C
		controller.ComicUpdate()
	}
}

func TaskChapter() {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()

	for {
		<-t.C
		wg := sync.WaitGroup{}
		wg.Add(10)
		for i := 0; i < 2; i++ {
			go func(i int) {
				controller.ChapterPaw()
				wg.Done()
			}(i)
		}
	}
}

func TaskImage() {
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()

	controller.ImagePaw()

	for {
		<-t.C
		controller.ImagePaw()
	}
}
