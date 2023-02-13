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
	Lock      sync.Mutex
	State     int
}

func SetUp() {
	num := config.Spe.Maxthreads
	if config.Spe.SourceId == 2 {
		num = num - 2
	}
	setRob(num)
}

func setRob(num int) {
	for {
		if len(Swarm) >= num {
			return
		}

		r := &Robot{
			Port: 19991 + len(Swarm),
		}
		r.prepare("https://" + config.Spe.SourceUrl)
		Swarm = append(Swarm, r)
	}
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

func ResetRob(Rob *Robot) {
	for tryLimit := 0; tryLimit <= 7; tryLimit++ {
		proxy := GetProxy()
		args := []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--ignore-certificate-errors",
			"--ignore-ssl-errors",
			"--user-agent=" + config.Spe.UserAgent,
		}
		if config.Spe.AppDebug == true {
			args = []string{
				"--ignore-certificate-errors",
				"--ignore-ssl-errors",
				"--user-agent=" + config.Spe.UserAgent,
			}
		}
		if proxy != "" {
			_ = append(args, proxy)
		}
		caps := selenium.Capabilities{
			"browserName": "chrome",
		}
		caps.AddChrome(chrome.Capabilities{
			Prefs: map[string]interface{}{
				"profile.managed_default_content_settings.images": 2,
			},
			Path: "",
			Args: args,
		})
		Rob.WebDriver.Close()
		wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Rob.Port))
		if err != nil {
			panic(err.Error())
		}
		err = wb.Get("https://" + config.Spe.SourceUrl)
		if err != nil {
			if tryLimit == 7 {
				wb.Close()
				panic(err.Error())
			}
		} else {
			Rob.WebDriver = wb
			Rob.WebDriver.ResizeWindow("", 1400, 1200)
			break
		}
	}
}

func (Robot *Robot) prepare(url string) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(config.Spe.SeleniumPath, Robot.Port, opts...)
	if nil != err {
		panic(err.Error())
	}
	Robot.Service = service

	for tryLimit := 0; tryLimit <= 7; tryLimit++ {
		proxy := GetProxy()
		args := []string{
			"--headless",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--ignore-certificate-errors",
			"--ignore-ssl-errors",
			"--user-agent=" + config.Spe.UserAgent,
		}
		if config.Spe.AppDebug == true {
			args = []string{
				"--ignore-certificate-errors",
				"--ignore-ssl-errors",
				"--user-agent=" + config.Spe.UserAgent,
			}
		}
		if config.Spe.AppDebug == false {
			args = append(args, proxy)
		}

		caps := selenium.Capabilities{
			"browserName": "chrome",
		}
		caps.AddChrome(chrome.Capabilities{
			Prefs: map[string]interface{}{
				"profile.managed_default_content_settings.images": 2,
			},
			Path: "",
			Args: args,
		})
		wb, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Robot.Port))
		if err != nil {
			panic(err.Error())
		}
		wb.SetImplicitWaitTimeout(time.Second * 10)
		wb.SetPageLoadTimeout(time.Second * 20)
		wb.ResizeWindow("", 1400, 1200)
		err = wb.Get(url)
		if err != nil {
			if tryLimit == 7 {
				wb.Close()
				service.Stop()
				panic(err.Error())
			}
		} else {
			Robot.WebDriver = wb
			break
		}
	}
}
