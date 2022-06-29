package contract

const (
	ConfigKeyLogLevel          = "log-level"
	ConfigKeyNoFail            = "no-fail"
	ConfigKeyNoFailAffirmative = "no-fail-affirmative"
	ConfigKeyOutput            = "output"
	ConfigKeyProbeTimeout      = "timeout"
)

type Config interface {
	String(key string, defaultValue string) string
	Bool(key string, defaultValue bool) bool
	Int(key string, defaultValue int) int
	IntSlice(key string, defaultValue []int) []int
	WithString(key, value string) Config
	WithBool(key string, value bool) Config
	WithInt(key string, value int) Config
	WithIntSlice(key string, value []int) Config
}
