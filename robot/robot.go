package robot

import (
	"comics/tools/config"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"sync"
	"time"
)

// selenium机器人池
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
	setting := func(lifeTime time.Time) {
		for {
			if len(Swarm) >= num {
				return
			}

			r := &Robot{
				Port:     19991 + len(Swarm),
				Lifetime: lifeTime,
			}
			r.Prepare("https://" + config.Spe.SourceUrl)
			Swarm = append(Swarm, r)
		}
	}
	lifeTime := time.Now().Add(time.Hour * 4)
	setting(lifeTime)

	t := time.NewTicker(time.Hour * 1)
	defer t.Stop()
	for {
		<-t.C
		activeNum := len(Swarm)
		if activeNum >= num {
			if Swarm[0].Lifetime.Second() < time.Now().Second() {
				continue //没有过期
			}
			for {
				if len(Swarm) == 0 {
					break
				}
				sw := pop(&Swarm)
				sw.Lock.Lock()
				sw.WebDriver.Close()
				sw.Service.Stop()
			}
		}
		setting(lifeTime)
	}
}

func pop(list *[]*Robot) *Robot {
	f := len(*list)
	rv := (*list)[f-1]
	*list = (*list)[:f-1]
	return rv
}

func (Robot *Robot) Prepare(url string) {
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
			"--user-agent=" + config.Spe.UserAgent,
		},
	}

	caps.AddChrome(chromeCaps)
	wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Robot.Port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = wb.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	Robot.WebDriver = wb
}
