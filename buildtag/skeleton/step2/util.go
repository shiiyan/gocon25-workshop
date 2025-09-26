//go:build go1.25

package main

import "fmt"

func init() {
	// TODO: go vetが検出するように%dを%sに変える
	fmt.Printf("wrong format %d", 100)
}
