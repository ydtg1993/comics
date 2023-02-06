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
