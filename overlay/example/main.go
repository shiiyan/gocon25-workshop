package main

import "time"

func main() {
	n := now()

	println(n.String())
}

func now() time.Time {
	return time.Now()
}
