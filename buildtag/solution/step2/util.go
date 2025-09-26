//go:build go1.25

package main

import "fmt"

func init() {
	// go vetが検出するはず
	fmt.Printf("wrong format %s", 100)
}
