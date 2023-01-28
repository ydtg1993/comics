package robot

import (
	"comics/tools/config"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Robot struct {
	Service   *selenium.Service
	WebDriver selenium.WebDriver
	Port      int
}

const SELENIUM_PATH = "chromedriver.exe"

func (Robot *Robot) Start(url string) {
	opts := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(SELENIUM_PATH, Robot.Port, opts...)
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
