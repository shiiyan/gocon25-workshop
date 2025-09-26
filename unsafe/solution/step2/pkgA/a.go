package pkgA

type A struct {
	s string // 追加された
	n int
}

func (a *A) N() int {
	return a.n
}
