package log

import "fmt"

func Alert(format string, args ...interface{}) string {
	return fmt.Sprintf("[ALERT] %s", fmt.Sprintf(format, args...))
}
