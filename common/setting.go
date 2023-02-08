package common

import (
	"comics/tools/config"
	"comics/tools/rd"
	"strconv"
	"time"
)

type Kv struct {
	Name string
	Val  int
}
type Kind struct {
	Tag    Kv
	Region Kv
	Pay    Kv
	State  Kv
}

const SourceChapterTASK = "source:comic:chapter"
const SourceChapterRetryTask = "source:comic:retry:chapter"

const SourceComicTASK = "source:comic:task"
const SourceComicRetryTask = "source:comic:retry:task"

const SourceImageTASK = "source:chapter:image"

const StopRobSignal = "stop"

func Signal(name string) bool {
	signal := rd.Get(StopRobSignal)
	if signal != "" {
		rd.Set(StopRobSignal, signal+"------"+strconv.Itoa(config.Spe.SourceId)+":"+name, time.Hour*2)
		return true
	} else {
		return false
	}
}
