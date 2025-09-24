// Go アセンブリで関数を実装します (ARM64版)
// 参考: https://go.dev/doc/asm

// func Add(a, b int64) int64
TEXT ·Add(SB), $0-24
	MOVD a+0(FP), R0    // 引数 a を R0 レジスタに読み込む
	MOVD b+8(FP), R1    // 引数 b を R1 レジスタに読み込む
	ADD R1, R0, R0      // R0 = R0 + R1
	MOVD R0, ret+16(FP) // 結果を戻り値の位置に書き込む
	RET

// func Sub(a, b int64) int64
TEXT ·Sub(SB), $0-24
	MOVD a+0(FP), R0    // 引数 a を R0 レジスタに読み込む
	MOVD b+8(FP), R1    // 引数 b を R1 レジスタに読み込む
	SUB R1, R0, R0      // R0 = R0 - R1 (引き算)
	MOVD R0, ret+16(FP) // 結果を戻り値の位置に書き込む
	RET
