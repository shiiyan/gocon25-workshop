package main

import (
	"fmt"
	"unsafe"

	"github.com/newmo-oss/gocon25-workshop/unsafe/solution/step1/pkgA"
)

func main() {

	type B struct {
		N int // 公開されたフィールド
	}

	var a pkgA.A
	// *A -> unsafe.Pointer -> *B
	b := (*B)(unsafe.Pointer(&a))
	b.N = 100

	// 100
	fmt.Println(a.N())
}
