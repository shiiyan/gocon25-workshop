author: Yuki Ito
summary: Delve 入門
id: delve
categories: codelab,markdown
environments: Web
status: Published

# Delve 入門

## Delve とは

[Delve](https://github.com/go-delve/delve) は Go 言語用のデバッガです。Delve を使うと、Go プログラムの実行を一時停止して変数の値を調べたり、プログラムをステップ実行したりすることができます。

## Delve をインストールしよう

```bash
> go install github.com/go-delve/delve/cmd/dlv@latest
```

インストール後にバージョンを確認してみましょう。

```bash
> dlv version
```

このワークショップで使う Delve のバージョンは `1.25.2` 以降となります。

## Delve を使ってみよう

まずは、簡単な Go プログラムをデバッグしてみましょう。

下記のコマンドでこのリポジトリに含まれるプログラムを実行してみてください。

```bash
> go run ./delve/example/calc
```

このプログラムは、2 つの整数値、x=2・y=3 をプログラム内で定義して、それぞれの値を 2 倍して足した結果を出力します。
出力結果は以下のようになります。

```text
13
```

おや...？計算結果はプログラムが出力している `13` ではなく `(2 * 2) + (3 * 2) = 10` のはずです。プログラムにバグがあるようです。

このプログラムを Delve でデバッグしてみましょう。
下記のコマンドを実行してみてください。

```bash
> dlv debug ./delve/example/calc
```

`dlv debug` コマンドは、指定した Go プログラムをコンパイルしてデバッグモードで実行します。
`dlv debug` コマンドでプログラムを実行すると、`(dlv)` というプロンプトが表示され、デバッグを行うための Delve のコマンドを入力できるようになります。

```text
Type 'help' for list of commands.
(dlv)
```

プログラムの実行を任意の場所で一時停止するために、下記のように `break` コマンドでブレークポイントを設定してみましょう。

```text
(dlv) break ./delve/example/calc/main.go:6
```

Delve の `break` コマンドはプログラムにブレークポイント（デバッグ中にプログラムを一時停止させたい位置）を設定するコマンドで、上記の例では `./delve/example/calc/main.go` の 6 行目にブレークポイントを設定しています（`break` コマンドは短縮した `b` でも実行できます）。

ブレークポイントを設定したら、下記のように `continue` コマンドでプログラムの実行を開始してみましょう。

```text
(dlv) continue
> [Breakpoint 1] main.main() ./delve/example/calc/main.go:6 (hits goroutine(1):1 total:1) (PC: ...)
     1: package main
     2:
     3: func main() {
     4:         var x, y int
     5:
=>   6:         x = 2
     7:         y = 3
     8:
     9:         answer := calc(x, y)
    10:
    11:         println(answer)
```

`continue` コマンドはプログラムの実行を開始するコマンドです（`continue` コマンドは短縮した `c` でも実行できます）。continue コマンドを実行すると、プログラムは最初のブレークポイントまで実行され、一時停止します。
上記の例では、先程設定したブレークポイントでプログラムが一時停止し、現在の実行位置が `./delve/example/calc/main.go` の 6 行目であることが表示されています。

このブレークポイントの位置では、変数 `x` に値 `2` を代入しているだけでバグではなさそうです。

次の処理に進んでみましょう。下記のように `step` コマンドを実行してみてください。

```text
(dlv) step
> main.main() ./delve/example/calc/main.go:7 (PC: ...)
     2:
     3: func main() {
     4:         var x, y int
     5:
     6:         x = 2
=>   7:         y = 3
     8:
     9:         answer := calc(x, y)
    10:
    11:         println(answer)
    12: }
```

`step` コマンドは現在の行を実行して次の処理に進むコマンドです（`step` コマンドは短縮した `s` でも実行できます）。上記の例では、現在の実行位置が `./delve/example/calc/main.go` の 7 行目に進んでいることが表示されています。

この行も変数 `y` に値 `3` を代入しているだけでバグではなさそうです。

さらに次の行に進んでみましょう。再度 `step` コマンドを実行してみてください。

```text
(dlv) step
> main.main() ./delve/example/calc/main.go:9 (PC: ...)
     4:         var x, y int
     5:
     6:         x = 2
     7:         y = 3
     8:
=>   9:         answer := calc(x, y)
    10:
    11:         println(answer)
    12: }
    13:
    14: func calc(x, y int) int {
```

この行では、`calc` 関数を呼び出して変数 `answer` に結果を代入しています。この `calc` 関数にバグがありそうです。

さらに `step` コマンドを実行して `calc` 関数の中に入ってみましょう。

```text
(dlv) step
> main.calc() ./delve/example/calc/main.go:14 (PC: ...)
     9:         answer := calc(x, y)
    10:
    11:         println(answer)
    12: }
    13:
=>  14: func calc(x, y int) int {
    15:         a := x + x
    16:         b := y * y
    17:
    18:         return a + b
    19: }
```

処理が `calc` 関数の中に入り、現在の実行位置が `./delve/example/calc/main.go` の 14 行目であることが表示されています。

この時点での引数 `x` と `y` の値を確認してみましょう。下記のように `print` コマンドで変数の値を表示してみてください。

```text
(dlv) print x
2
(dlv) print y
3
```

`print` コマンドは指定した変数の値を表示するコマンドです（`print` コマンドは短縮した `p` でも実行できます）。上記の例では、引数 `x` の値が `2`、引数 `y` の値が `3` であることが表示されています。

それでは、`calc` 関数の中をさらにステップ実行してみましょう。`step` コマンドを 2 回実行してみてください。

```text
(dlv) step
> main.calc() ./delve/example/calc/main.go:15 (PC: ...)
    10:
    11:         println(answer)
    12: }
    13:
    14: func calc(x, y int) int {
=>  15:         a := x + x
    16:         b := y * y
    17:
    18:         return a + b
    19: }

(dlv) step
> main.calc() ./delve/example/calc/main.go:16 (PC: ...)
    11:         println(answer)
    12: }
    13:
    14: func calc(x, y int) int {
    15:         a := x + x
=>  16:         b := y * y
    17:
    18:         return a + b
    19: }
```

変数 `a` への代入が実行されたので、変数 `a` の値を確認してみましょう。再度 `print` コマンドで変数 `a` の値を表示してみてください。

```text
(dlv) print a
4
```

変数 `a` の値が `4` であることが表示されました。x=2 の 2 倍は 4 なので、ここまでは正しい計算が行われています。

それでは、次に変数 `b` への代入を実行してみましょう。再度 `step` コマンドを実行してみてください。

```text
(dlv) step
> main.calc() ./delve/example/calc/main.go:18 (PC: ...)
    13:
    14: func calc(x, y int) int {
    15:         a := x + x
    16:         b := y * y
    17:
=>  18:         return a + b
```

変数 `b` への代入が実行され、現在の実行位置が `./delve/example/calc/main.go` の 18 行目に進んでいることが表示されています。

それでは、変数 `b` の値を確認してみましょう。再度 `print` コマンドで変数 `b` の値を表示してみてください。

```text
(dlv) print b
9
```

変数 `b` の値が `9` であることが表示されました。y=3 の 2 倍は 6 なので、ここでバグが発生していることがわかります。

現在の実行位置のソースコードを再表示するために、`list` コマンドを実行してみましょう。

```text
(dlv) list
> main.calc() ./delve/example/calc/main.go:18 (PC: 0x477b0c)
    13:
    14: func calc(x, y int) int {
    15:         a := x + x
    16:         b := y * y
    17:
=>  18:         return a + b
    19: }
```

`list` コマンドは現在の実行位置のソースコードを表示するコマンドです（`list` コマンドは短縮した `l` でも実行できます）。上記の例では、現在の実行位置が `./delve/example/calc/main.go` の 18 行目であることが表示されています。

ここで、変数 `b` の値を計算している 16 行目のソースコードを確認してみると、計算式が `y * y` となっており、y の 2 乗を計算していることがわかります。正しくは y の 2 倍を計算するために `y + y` とする必要があり、この部分がバグであることがわかりました！

最後に、デバッグが完了したので Delve を終了しましょう。`quit` コマンドを実行してみてください。

```text
(dlv) quit
```

`quit` コマンドを実行すると、Delve が終了します（`quit` コマンドは短縮した `q` でも実行できます）。

このように、Delve デバッガを使うと Go プログラムの実行を一時停止して変数の値を調べたり、プログラムをステップ実行したりすることができるので、プログラムのバグの特定に役立ちます。

## テストでの Delve の利用

先ほどの例では、Delve を使って main 関数からはじまるプログラム全体のデバッグを行いましたが、実践的には main 関数を持たないライブラリのデバッグを行いたい場合も多いでしょう。

そのような場合は、テストコードでライブラリの機能を呼び出して、テスト経由で Delve を利用する方法が便利です。

前のステップでデバッグしていた `calc` 関数のためのテストコードは `./delve/example/calc/calc_test.go` に以下のように実装されています。

```go
package main

import "testing"

func TestCalc(t *testing.T) {
	t.Parallel()

	want := 10
	got := calc(2, 3)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
```

このテストコードを利用して `calc` 関数をデバッグしてみましょう。下記の コマンドを実行してみてください。

```bash
> dlv test ./delve/example/calc
Type 'help' for list of commands.                                                                                                                                                                                 (dlv)
```

前のステップで利用した `dlv debug` コマンドと同じように、`dlv test` コマンドは指定したパッケージのテストに対してデバッガを起動します。

デバッガが起動されたあとは、前のステップと同じように Delve の各種コマンドを利用してデバッグを行うことができます。

```
(dlv) break ./delve/example/calc/main.go:15
Breakpoint 1 set at 0x5a7ada for github.com/newmo-oss/gocon25-workshop/delve/example/calc.calc() ./delve/example/calc/main.go:15

(dlv) c
> [Breakpoint 1] github.com/newmo-oss/gocon25-workshop/delve/example/calc.calc() ./delve/example/calc/main.go:15 (hits goroutine(21):1 total:1) (PC: ...)
    10:
    11:         println(answer)
    12: }
    13:
    14: func calc(x, y int) int {
=>  15:         a := x + x
    16:         b := y * y
    17:
    18:         return a + b
    19: }

(dlv) p x
2
```

また、`dlv test` コマンドでは `--` を指定して Go のテスト（`go test`）に渡すオプションを指定することもできます。

実践的には、テストとデバッグの実行時間を短くするために、特定のテスト関数だけを実行したい場合が多いでしょう。そのような場合、下記のように `--` を利用して `-test.run` オプションで実行したいテスト関数を指定することができます。

```bash
> dlv test ./delve/example/calc -- -test.run TestCalc
```

このように、`dlv test` コマンドを利用すると、main 関数を経由することなくテストを用いて Go のプログラムをデバッグすることができます。

## まとめ

このワークショップでは、Go 言語用のデバッガである Delve のインストール方法と基本的な使い方について学びました。

Delve デバッガを利用すると、Go のプログラムをステップ実行したり、実行中の変数の値を調べたりすることができるので、プログラムのバグをより素早く特定できるようになります。

特に、 `dlv test` コマンドは特定の関数やメソッドをテストを通じてデバッグすることができるので、実践的な Go の開発において非常に役立ちます。

Delve は今回紹介したもの以外にも多くの便利な機能を持っています。ぜひ[公式のドキュメント](https://github.com/go-delve/delve/tree/master/Documentation)や `help` コマンドを参照して、Delve のさらなる活用方法を学んでみてください！
