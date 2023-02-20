package common

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

/*redis key*/
const SourceChapterTASK = "source:comic:chapter"
const SourceChapterRetryTask = "source:comic:retry:chapter"

const SourceComicTASK = "source:comic:task"
const SourceComicRetryTask = "source:comic:retry:task"

const SourceImageTASK = "source:chapter:image"

const SourceImageCapture = "source:image:capture"
const SourceImageDownload = "source:image:download"

const TaskStepRecord = "task:step:record"

const StopRobotSignal = "shutdown"
