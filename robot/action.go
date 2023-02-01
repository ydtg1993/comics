package robot

import (
	"time"
)

func (Robot *Robot) TapIn(method, path string) bool {
	for i := 0; i < 3; i++ {
		dom, err := Robot.WebDriver.FindElement(method, path)
		if err != nil {
			continue
		}
		dom.Click()
		time.Sleep(1 * time.Second)
		//检查dom 判断跳转
		_, err = Robot.WebDriver.FindElement(method, path)
		if err == nil {
			return true
		}
	}
	return false
}
