package gdiscovery

type IConfig interface {
	String(path string) string
	MustString(path string) string
	Int64(path string) int64
	Bool(path string) bool
}
