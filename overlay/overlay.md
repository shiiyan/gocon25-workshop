author: Yuki Ito
summary: overlay 入門
id: overlay
categories: codelab,markdown
environments: Web
status: Published

# overlay 入門

## overlay とは

Go のビルドでは、overlay オプションを使うことで特定のファイルを別のファイルで置き換えてビルドすることができます。

このワークショップでは、overlay の基本的な使い方とテストでの実践的な利用方法について学びます。

## overlay をつかってみよう

まずは、簡単な Go プログラムを題材に overlay を使ってみましょう。

下記のコマンドでこのリポジトリに含まれるプログラムを実行してみてください。

```text
> go run ./overlay/example
```

このプログラムは、現在時間を `time.Now()` で取得して、それを表示するだけのシンプルなプログラムです。

```go
package main

import "time"

func main() {
	n := now()

	println(n.String())
}

func now() time.Time {
	return time.Now()
}
```

表示される時間はプログラムを実行した現在の時間となるので、毎回異なる値が表示されます。

この現在時刻について、「プログラムの実装に変更を入れることなく実験などのために現在時刻を固定したい」ということを想定してみましょう。

このような場合に、Go のビルドの overlay 機能を使うことで標準パッケージの `time.Now` を置き換えて固定の時刻を返すようにすることができます。

まずは、overlay のための設定ファイル（JSON）を見てみましょう。

```json
> cat ./overlay/example/overlay.json
{
  "Replace":{
    "${GOROOT}/src/time/time.go": "./overlay/example/time/time.go"
  }
}
```

※ `${GOROOT}` は Go のインストールディレクトリ（`go env GOROOT` コマンドで取得できる値）に置き換えてください。

この JSON ファイルは、`Replace` というキーを持つオブジェクトを定義しています。このオブジェクトの中で、`time/time.go` という標準パッケージのファイルを、このリポジトリに含まれる `./overlay/example/time/time.go` という別のファイルで置き換えることを指定しています。

この設定ファイルを使って overlay を利用してみましょう。下記のコマンドを実行してみてください。

```text
> go run -overlay ./overlay/example/overlay.json ./overlay/example
2025-09-27 00:00:00 +0000 UTC
```

このコマンドでは先程と同じプログラムを実行していますが、表示される時間が `2025-09-27 00:00:00 +0000 UTC` に固定されていることがわかります。異なる部分は、`-overlay` オプションで `overlay.json` ファイルを指定していることです。

`go build` や `go test` では、`-overlay` というオプションを提供しています。このオプションは、指定した JSON ファイルの内容に基づいてビルド時に特定のファイルを別のファイルで置き換えるためのオプションです。

このリポジトリに含まれる `./overlay/example/time/time.go` ファイルは、標準パッケージに含まれる `time/time.go` ファイルをコピーして、下記のように `Now` 関数を固定の時刻を返すように書き換えたものです。

```go
// ...

func Now() Time {
	return Date(2025, 9, 27, 0, 0, 0, 0, UTC)
}

// ...
```

`-overlay` オプションを使うことで、この書き換えた `Now` 関数がビルド時に利用されるようになり、プログラムの実行時に `time.Now()` を呼び出すと固定の時刻が返されるようになります。

このように、Go のビルドが提供している overlay 機能を使うことで特定のファイルを別のファイルで置き換えてビルドすることができるので、プログラムの実装に変更を加えることなく動作を変更することができます。

## テストでの overlay の利用

先程の例では、`go run` コマンドで main 関数からはじまるプログラムの実行に overlay を利用してソースコードの一部を置き換えました。

より実践的な利用例として、テストコードで overlay を利用する方法を見てみましょう。下記のテストコマンドを実行してみてください。

```text
> go test ./overlay/example

--- FAIL: TestNow (0.00s)
...
```

このテストは先程見ていた `now` 関数のためのテストで、下記のように実装されています。

```go
package main

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	t.Parallel()

	want := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	got := now()

	if !got.Equal(want) {
		t.Errorf("now() = %v; want %v", got, want)
	}
}
```

このテストコードは、`now` 関数が返す現在時刻が `2025-09-27 00:00:00 +0000 UTC` であることを期待していますが、通常のビルドでは `now` 関数は `time.Now()` を呼び出して現在の時刻を返すため、テストは失敗します。

このテストを overlay を使って実行してみましょう。下記のコマンドを実行してみてください。

```text
> go test -overlay ./overlay/example/overlay.json ./overlay/example
ok      github.com/newmo-oss/gocon25-workshop/overlay/example   0.001s
```

この `go test` コマンドでは、先程のステップと同じように `-overlay` オプションを指定しているため、`time.Now()` の実装が置き換えられて固定の時刻を返すようになり、テストが成功していることがわかります。

このように、テストにおいて overlay を利用することで実装を変更することなく特定の関数の動作を置き換えてテストを実行することができるようになり、場合によっては通常のビルドでは難しいテストを実行できるようになります。

## まとめ

このワークショップでは、Go のビルドが提供している overlay 機能の基本的な使い方と、テストでの実践的な利用方法について学びました。

日常的な開発において overlay 機能を多用することはあまりないかもしれませんが、特定の関数の動作を置き換えたい場合などに役立つことがあるので一つの解決策として覚えておくと良いでしょう。
