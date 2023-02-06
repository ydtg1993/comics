package robot

import (
	"comics/tools"
	"comics/tools/config"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"sync"
	"time"
)

var Swarm []*Robot

type Robot struct {
	Service   *selenium.Service
	WebDriver selenium.WebDriver
	Port      int
	Lifetime  time.Time
	Lock      sync.Mutex
	State     int
}

func SetUp(num int) {
	lifeTime := time.Now().Add(time.Hour * 6)
	setRob(num, lifeTime)

	t := time.NewTicker(time.Minute * 30)
	defer t.Stop()
	for {
		<-t.C
		activeNum := len(Swarm)
		if activeNum >= num {
			if Swarm[0].Lifetime.Second() < time.Now().Second() {
				continue //没有过期
			}
			deleteRob()
		}
		setRob(num, lifeTime)
	}
}

func ReSetUp(num int) {
	deleteRob()
	lifeTime := time.Now().Add(time.Hour * 6)
	setRob(num, lifeTime)
}

func GetRob(keys []int) *Robot {
	var Rob *Robot
	for k, robot := range Swarm {
		if len(keys) > 0 {
			exists, _ := tools.InArray(k, keys)
			if exists == false {
				continue
			}
		}
		if robot.State == 1 {
			continue
		}
		robot.Lock.Lock()
		robot.State = 1
		Rob = robot
		break
	}
	return Rob
}

func ResetRob(rob *Robot) {
	rob.State = 0
	rob.Lock.Unlock()
}

func setRob(num int, lifeTime time.Time) {
	for {
		if len(Swarm) >= num {
			return
		}

		r := &Robot{
			Port:     19991 + len(Swarm),
			Lifetime: lifeTime,
		}
		r.prepare("https://" + config.Spe.SourceUrl)
		Swarm = append(Swarm, r)
	}
}

func deleteRob() {
	for {
		if len(Swarm) == 0 {
			return
		}
		sw := pop(&Swarm)
		sw.WebDriver.Close()
		sw.Service.Stop()
		sw.State = 1
	}
}

func (Robot *Robot) prepare(url string) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(config.Spe.SeleniumPath, Robot.Port, opts...)
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	Robot.Service = service
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			//"--headless",
			//"--no-sandbox",
			"--ignore-certificate-errors",
			"--ignore-ssl-errors",
			"--user-agent=" + config.Spe.UserAgent,
		},
	}

	caps.AddChrome(chromeCaps)
	wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Robot.Port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	wb.ResizeWindow("", 1400, 900)
	err = wb.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Robot.WebDriver = wb
}

func pop(list *[]*Robot) *Robot {
	f := len(*list)
	rv := (*list)[f-1]
	*list = (*list)[:f-1]
	return rv
}
