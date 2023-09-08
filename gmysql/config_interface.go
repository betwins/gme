package gmysql

type IConfig interface {
	String(path string) string
	MustString(path string) string
	Int64(path string) int64
	Int(path string) int
	Bool(path string) bool
}
