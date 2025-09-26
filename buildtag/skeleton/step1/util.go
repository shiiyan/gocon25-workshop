//TODO: go1.25以上を対象としたビルドタグを記述する

package main

import (
	"fmt"
	"runtime"
)

func init() {
	msg := fmt.Sprintf("panic in %s", runtime.Version())
	panic(msg)
}
