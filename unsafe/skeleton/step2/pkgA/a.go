package pkgA

type A struct {
	// TODO: string型のフィールドsを追加する
	n int
}

func (a *A) N() int {
	return a.n
}
