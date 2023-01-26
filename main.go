package main

import (
	"comics/controller"
	"comics/robot"
	"github.com/tebeka/selenium"
)

func main() {
	controller.ComicPaw()

	done := make(chan int)
	Robot := new(robot.Robot)
	defer Robot.Service.Stop()
	Robot.Start("https://www.youtube.com/")
	Robot.TapIn(selenium.ByXPATH, "//*[@id='HomeContributeArea']/div[2]/div[1]/a/div[1]")
	<-done
}
