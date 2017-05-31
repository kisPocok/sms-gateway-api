package logger

import "fmt"

// Logging define a basic logger class
type Logging interface {
	Log(s ...string)
}

// Logger represents a simple reporter
type Logger struct{}

// Log echoes msg
func (l Logger) Log(s ...string) {
	fmt.Println(s)
}

// Null logs nothing
type Null struct{}

// Log echoes msg
func (l Null) Log(s ...string) {}
