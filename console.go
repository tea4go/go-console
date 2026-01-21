package console

import (
	"github.com/tea4go/go-console/interfaces"
)

// Console communication interface
type Console interfaces.Console

// New creates a new console with initial size
func New(w int, h int) (Console, error) {
	return newNative(w, h)
}
