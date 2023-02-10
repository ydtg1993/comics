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

const SourceChapterTASK = "source:comic:chapter"
const SourceChapterRetryTask = "source:comic:retry:chapter"

const SourceComicTASK = "source:comic:task"
const SourceComicRetryTask = "source:comic:retry:task"

const SourceImageTASK = "source:chapter:image"
