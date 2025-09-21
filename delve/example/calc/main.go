package main

func main() {
	var x, y int

	x = 2
	y = 3

	answer := calc(x, y)

	println(answer)
}

func calc(x, y int) int {
	a := x + x
	b := y * y

	return a + b
}
