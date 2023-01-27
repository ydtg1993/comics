package main

import (
	"comics/controller"
	_ "comics/rd"
	"comics/tools/config"
	"comics/tools/database"
	"comics/tools/log"
	"comics/tools/rd"
	"github.com/beego/beego/v2/core/logs"
	"runtime"
	"time"
)

func main() {
	Setup()

	TaskComic()

	//go TaskChapter()

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

	//启动更新数据到es
	controller.ComicPaw()

	for {
		<-t.C
		controller.ComicPaw()
	}
}

func TaskChapter() {
	t := time.NewTicker(time.Minute * 5)
	defer t.Stop()

	controller.ChapterPaw()

	for {
		<-t.C
		controller.ChapterPaw()
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
