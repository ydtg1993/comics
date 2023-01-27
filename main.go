package main

import (
	"comics/controller"
	_ "comics/log"
	_ "comics/rd"
)

func main() {
	controller.ComicPaw()

	/*done := make(chan int)
	Robot := new(robot.Robot)
	defer Robot.Service.Stop()
	Robot.Start("https://ac.qq.com/Comic/all/page/1")
	Robot.TapIn(selenium.ByXPATH, "/html/body/div[3]/div[2]/div/div[2]/ul/li[1]/div[1]/a")
	<-done*/
}
