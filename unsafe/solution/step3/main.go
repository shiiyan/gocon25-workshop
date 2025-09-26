package main

import (
	"fmt"
	"unsafe"
)

func main() {
	array := [...]int{10, 20, 30, 40, 50}
	slice := array[:]

	type SliceHeader struct {
		Ptr unsafe.Pointer
		Len int
		Cap int
	}

	header := (*SliceHeader)(unsafe.Pointer(&slice))
	// main.SliceHeader{Ptr:(unsafe.Pointer)(0xXXXXXXXXXX), Len:5, Cap:5}
	fmt.Printf("%#v\n", *header)

	ptr := unsafe.Add(header.Ptr, unsafe.Sizeof(0)*2)
	header2 := &SliceHeader{
		Ptr: ptr,
		Len: 2,
		Cap: 3,
	}
	// *SliceHeader -> []int
	slice2 := *(*[]int)(unsafe.Pointer(header2))
	// [30 40] 2 3
	fmt.Println(slice2, len(slice2), cap(slice2))
}
