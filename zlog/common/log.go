package common

// OneLog a log
type OneLog struct {
	OutTo      int
	FormatType int
	Level      int
	CallerFile string
	CallerLine int
	CallerName string
	CallerPkg  string
	Timestamp  int64
	Format     string
	Args       []interface{}
}
