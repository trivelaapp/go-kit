package log

// Level indicates the severity of the data being logged.
type Level int

const (
	// LevelCritical alerts about severe problems. Most of the time, needs some human intervention ASAP.
	LevelCritical Level = iota + 1
	// LevelError alerts about events that are likely to cause problems.
	LevelError
	// LevelWarning warns about events the might cause problems to the system.
	LevelWarning
	// LevelInfo are routine information.
	LevelInfo
	// LevelDebug are debug or trace information.
	LevelDebug
)

var levelStringValueMap = map[string]Level{
	"CRITICAL": LevelCritical,
	"ERROR":    LevelError,
	"WARNING":  LevelWarning,
	"INFO":     LevelInfo,
	"DEBUG":    LevelDebug,
}

// String returns the name of the LogLevel.
func (l Level) String() string {
	return []string{
		"CRITICAL",
		"ERROR",
		"WARNING",
		"INFO",
		"DEBUG",
	}[l-1]
}
