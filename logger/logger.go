package logger

import "fmt"

// Logging define a basic logger class
type Logging interface {
	Log(s string)
}

// Logger represents a simple reporter
type Logger struct {
	Disabled bool
}

// Log echoes msg
func (l Logger) Log(s string) {
	if !l.Disabled {
		fmt.Println(s)
	}
}
