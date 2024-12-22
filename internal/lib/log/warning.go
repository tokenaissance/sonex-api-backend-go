package log

import "fmt"

func Warning(format string, args ...interface{}) string {
	return fmt.Sprintf("[WARNING] %s", fmt.Sprintf(format, args...))
}
