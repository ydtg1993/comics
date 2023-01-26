package log

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"runtime"
	"strings"
)

func init() {

	//日志天数
	var logDay = ""
	//日志路径
	var logpath = ""
	//日志级别
	var loglevel = ""
	logErr := "error"

	if logDay = os.Getenv("logday"); logDay == "" {
		logDay = "7"
	}
	if logpath = os.Getenv("logpath"); logpath == "." {
		logpath = ""
	}
	//日志级别
	loglevel = os.Getenv("loglevel")

	logName := fmt.Sprintf("%shdb-import.log", logpath)

	if len(loglevel) > 1 {
		sLevel := strings.Split(loglevel, ",")
		logErr = `"` + strings.Join(sLevel, `","`) + `"`
	}

	level := 2
	if strings.ContainsAny(loglevel, "debug") == true {
		level = 7
	}

	logCfg := fmt.Sprintf(`{"filename":"%s","level":%d,"maxdays":%s,"separate":[%s]}`, logName, level, logDay, logErr)
	//记录日志
	err := logs.SetLogger(logs.AdapterMultiFile, logCfg)
	if err != nil {
		panic(err)
	}

	// 开始前的线程数
	logs.Debug("线程数量 starting: %d\n", runtime.NumGoroutine())
}
