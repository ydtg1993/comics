package log

import (
	"comics/tools/config"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"strings"
)

type LogsManage struct {
}

func (l LogsManage) SetUp() (err error) {
	logName := fmt.Sprintf("%shdb-import.log", config.Spe.Logpath)
	logErr := `"` + strings.Join(config.Spe.Loglevel, `","`) + `"`
	level := 2
	if strings.ContainsAny(logErr, "debug") == true {
		level = 7
	}
	logCfg := fmt.Sprintf(`{"filename":"%s","level":%d,"maxdays":%d,"separate":[%s]}`, logName, level, config.Spe.Logday, logErr)
	//记录日志
	err = logs.SetLogger(logs.AdapterMultiFile, logCfg)
	return
}
