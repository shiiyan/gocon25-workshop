package pkgA

type A struct {
	n int
}

func (a *A) N() int {
	return a.n
}
