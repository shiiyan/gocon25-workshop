//go:build go1.25

package main

import (
	"fmt"
	"runtime"
)

func init() {
	msg := fmt.Sprintf("panic in %s", runtime.Version())
	panic(msg)
}
