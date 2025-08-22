package util

import (
	"fmt"
	"time"
)

// MakeLog creates a log message with a timestamp.
func MakeLog(message string, args ...any) string {
	return fmt.Sprintf("[%s] %s", time.Now().Local().String(), fmt.Sprintf(message, args...))
}
