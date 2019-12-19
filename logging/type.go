package logging

// LogType is an enum for different log levels
type LogType int

const (
	// Debug is the LogType for base messages for debugging purposes
	Debug LogType = 0
	// Info is the LogType for normal messages
	Info LogType = 1
	// Failure is the LogType for failed messages
	Failure LogType = 2
)

// SetLogLevel sets the global LogLevel for logging
func SetLogLevel(value int) {
	switch value {
	case 0:
		logLevel = Debug
	case 1:
		logLevel = Info
	case 2:
		logLevel = Failure
	default:
		logLevel = Debug
	}
}
