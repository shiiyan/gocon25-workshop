author: tenntenn
summary: Introduction to buildtag bomb
id: buildtag
categories: codelab,markdown
environments: Web
status: Published
url: gocon25-workshop

# ビルドタグを使った時限爆弾

このコードラボでは、ビルドタグを使った攻撃方法とその対応方法について学習します。

## 学習目標

- ビルドタグについて学ぶ
- ビルドタグを使った攻撃方法について学ぶ
- ビルドタグを使った攻撃への対応方法について学ぶ

## 進め方

1. 各ステップの「学習内容」を読み、概念を理解します
2. `skeleton/stepX/*.go` のTODOコメントを修正して実装します
3. 修正したコードを実行します
4. `solution/stepX/*.go` と比較して理解を深めます

### 前提条件

- 基本的な Go の文法の理解

### 実行して確認

各ステップのコードを修正したら、以下のコマンドで動作を確認してください。
`GOTOOLCHAIN`環境変数で指定している`Go`のバージョンは各ステップの指示に従ってください。

```bash
GOTOOLCHAIN=go1.xx.yy go run ./skeleton/stepX
```

---

## Step 1: ビルドタグと攻撃手法

### このステップで学ぶこと

このステップでは、`Go`のビルドタグ（build constraint）を悪用した時限爆弾攻撃の仕組みを学びます。将来の`Go`バージョンで実行される悪意のあるコードを仕込む手法と、それが何故危険なのかを理解します。

### Build constraint（ビルドタグ）

ビルドタグは、特定の条件下でのみファイルをビルドに含めるための仕組みです。

```go
//go:build go1.25

package main

func init() {
    panic("This code runs only with Go 1.25 or later")
}
```

主な用途として以下があります。

- プラットフォーム固有のコード（`//go:build windows`）
- Goバージョン固有のコード（`//go:build go1.18`）
- 開発・本番環境の切り替え（`//go:build debug`）

### `GOTOOLCHAIN`環境変数

`GOTOOLCHAIN`環境変数を使用すると、特定のGoバージョンでコードを実行できます。

```bash
$ GOTOOLCHAIN=go1.24.6 go run .   # Go 1.24.6で実行
$ GOTOOLCHAIN=go1.25.1 go run .   # Go 1.25.1で実行
```

`Go 1.21`以降、この機能により必要に応じて新しい`Go`バージョンが自動的にダウンロードされます。

### ビルドタグを使った攻撃手法

悪意のある開発者は、将来の`Go`バージョンでのみ実行される悪質なコードを仕込むことができます。

```go
//go:build go1.25

package main

func init() {
    // 将来のバージョンでのみ実行される悪意のあるコード
    panic("Time bomb activated!")
}
```

この攻撃の特徴は次の通りです。

1. **遅延実行** - 現在は無害に見えるが、将来の`Go`バージョンで悪意のあるコードが実行される
2. **検出困難** - 現在の`Go`バージョンでは該当ファイルがビルドされないため、テストで発見しにくい
3. **広範囲影響** - 依存関係を通じて多くのプロジェクトに影響を与える可能性

**重要な概念**

- ビルドタグは条件付きコンパイルを可能にする強力な機能
- `GOTOOLCHAIN`環境変数により任意のGoバージョンで実行可能
- 将来のバージョンを指定したビルドタグは時限爆弾として悪用される可能性
- 依存関係のセキュリティ監査には将来のビルドタグの確認も必要

### 実装タスク

このステップでは、ビルドタグを使った時限爆弾を作成し、その動作を確認します。

`skeleton/step1/util.go`を修正して、次の要件を満たしてください。

1. Go 1.25以上でのみ実行されるビルドタグを追加
2. 異なるGoバージョンで実行して動作の違いを確認

実行例：
```bash
# Go 1.24では正常実行（util.goが読み込まれない）
GOTOOLCHAIN=go1.24.6 go run ./skeleton/step1

# Go 1.25では panic が発生（util.goが読み込まれる）
GOTOOLCHAIN=go1.25.1 go run ./skeleton/step1
```

### 理解度チェック

- ビルドタグがどのような仕組みで動作するか説明できますか？
- なぜ将来のGoバージョンを指定したビルドタグが危険なのか理解していますか？
- この攻撃手法を防ぐためにはどのような対策が考えられますか？
- `GOTOOLCHAIN`環境変数の役割と使い方を説明できますか？

---

## Step 2: 静的解析とビルドタグ

### このステップで学ぶこと

このステップでは、ビルドタグが静的解析ツールに与える影響を学びます。ビルドタグによりファイルが条件付きでビルドされるため、静的解析も非対称になり、潜在的な問題を見逃す可能性があることを理解します。

### 静的解析が対象とするファイル

静的解析ツール（`go vet`、[govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)、[Google Capslock](https://github.com/google/capslock)など）は、現在のビルド環境でコンパイル対象となるファイルのみを解析します。

```bash
# Go 1.24では util.go は解析対象外
$ GOTOOLCHAIN=go1.24.6 go vet ./skeleton/step2
# エラーが検出されない

# Go 1.25では util.go も解析対象
$ GOTOOLCHAIN=go1.25.1 go vet ./skeleton/step2
# エラーが検出される
```

これにより、将来のバージョンでのみ実行される脆弱なコードや問題のあるコードが静的解析で見逃される危険性があります。

### ビルドタグの検出とチェック

ビルドタグ付きファイルを検出するには、専用のツールが必要です。

#### gostaticanalysis/buildtagツール

[gostaticanalysis/buildtag](https://github.com/gostaticanalysis/buildtag)は、ビルドタグ付きファイルを検出・解析するツールです。

```bash
$ go install github.com/gostaticanalysis/buildtag/cmd/buildtag@latest
$ buildtag ./...
```

#### go/build/constraintパッケージ

Go標準ライブラリの`go/build/constraint`パッケージを使用してビルドタグを解析できます。

```go
import "go/build/constraint"

// ビルドタグの解析例
expr, err := constraint.Parse("//go:build go1.25")
if err != nil {
    log.Fatal(err)
}
// expr.Eval(map[string]bool{"go1.25": true}) // true
```

**重要な概念**

- 静的解析ツールは現在のビルド環境でのみファイルを検査する
- ビルドタグ付きファイルは条件によって解析対象から除外される
- セキュリティツールも同様の制限を受ける
- 包括的なセキュリティ監査にはすべてのビルドタグ組み合わせでの検査が必要

### 実装タスク

このステップでは、ビルドタグが静的解析に与える影響を確認します。

`skeleton/step2/util.go`を修正して、次の要件を満たしてください。

1. `fmt.Printf`のフォーマット文字列を意図的に間違える（`%d`を`%s`に変更）
2. この変更により`go vet`の`printf`アナライザーがエラーを検出するようにする

実行例：
```bash
# Go 1.24では util.go が解析されないためエラーなし
$ GOTOOLCHAIN=go1.24.6 go vet ./skeleton/step2

# Go 1.25では util.go が解析されエラーが検出される
$ GOTOOLCHAIN=go1.25.1 go vet ./skeleton/step2
./skeleton/step2/util.go:9:2: Printf format %s has arg 100 of wrong type int
```

### 理解度チェック

- なぜ静的解析ツールはビルドタグ付きファイルを見逃すのか説明できますか？
- セキュリティツールが同様の制限を受ける理由を理解していますか？
- ビルドタグ付きファイルを包括的にチェックするにはどうすれば良いですか？
- `go/build/constraint`パッケージの役割を説明できますか？

---


## まとめ

### 学んだ概念

#### Step 1: ビルドタグと攻撃手法

- ビルドタグ（build constraint）は条件付きコンパイルを可能にする機能
- `GOTOOLCHAIN`環境変数により任意のGoバージョンでの実行が可能
- 将来のGoバージョンを指定したビルドタグによる時限爆弾攻撃の仕組み
- 攻撃の特徴：遅延実行、検出困難、広範囲影響

#### Step 2: 静的解析とビルドタグ

- 静的解析ツールは現在のビルド環境でのみファイルを解析する
- ビルドタグ付きファイルは条件によって解析対象から除外される
- `go vet`コマンドも同様の制限を受ける
- セキュリティツールによる包括的な監査の必要性

### 次のステップ

1. 実際のOSSプロジェクトでビルドタグ付きファイルの存在を確認する
2. `gostaticanalysis/buildtag`ツールを使ってビルドタグの検出を試す
3. CI/CDパイプラインで複数のGoバージョンでの静的解析を実装する
4. 依存関係のセキュリティ監査にビルドタグチェックを組み込む

## 参考資料

- [Go Build Constraints](https://pkg.go.dev/go/build#hdr-Build_Constraints)
- [go/build/constraint package](https://pkg.go.dev/go/build/constraint)
- [gostaticanalysis/buildtag](https://github.com/gostaticanalysis/buildtag)
- [GOTOOLCHAIN環境変数](https://go.dev/doc/toolchain)
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
- [Google Capslock](https://github.com/google/capslock)
