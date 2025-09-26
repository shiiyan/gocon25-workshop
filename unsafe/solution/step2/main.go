package main

import (
	"fmt"
	"unsafe"

	"github.com/newmo-oss/gocon25-workshop/unsafe/solution/step2/pkgA"
)

func main() {

	type B struct {
		N int // 公開されたフィールド
	}

	// Aのサイズが変わると-delta*deltaが0じゃなくなりコンパイルエラーになる
	const delta = int64(unsafe.Sizeof(B{})) - int64(unsafe.Sizeof(pkgA.A{}))
	var _ [-delta * delta]int

	var a pkgA.A
	// *A -> unsafe.Pointer -> *B
	b := (*B)(unsafe.Pointer(&a))
	b.N = 100

	// 100
	fmt.Println(a.N())
}
