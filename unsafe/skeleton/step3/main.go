package main

import (
	"fmt"
	"unsafe"
)

func main() {
	array := [...]int{10, 20, 30, 40, 50}
	slice := array[:]

	type SliceHeader struct {
		Ptr // TODO: 配列のポインタを保持するための型
		Len int
		Cap int
	}

	header := (/* TODO: 変換したい型 */)(unsafe.Pointer(&slice))
	// main.SliceHeader{Ptr:(unsafe.Pointer)(0xXXXXXXXXXX), Len:5, Cap:5}
	fmt.Printf("%#v\n", *header)

	// TODO: go vetでエラーが出ないようにunsafe.Addを使用する
	ptr := uintptr(header.Ptr) + unsafe.Sizeof(0)*2
	header2 := &SliceHeader{
		Ptr: unsafe.Pointer(ptr), // TODO: 不要な型変換を取り除く
		Len: 2,
		Cap: 3,
	}
	// *SliceHeader -> []int
	slice2 := *(*[]int)(unsafe.Pointer(header2))
	// [30 40] 2 3
	fmt.Println(slice2, len(slice2), cap(slice2))
}
