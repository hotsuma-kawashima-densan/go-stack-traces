# エラー出力の実装方法について

## bad

下記のようにService層とServer層の両方でエラー出力が実装されている場合、同じエラーが二重で出力されてしまいます。

##### handler

```go
func handler() {
	_, err := service()

	if err != nil {
        // ↓ログを記録し、エラーを返す
		slog.Error("エラーが発生しました", "cause", err)
	}
}
```

##### service

```go
func service() (interface{}, error) {
	err := query()

	if err != nil {
        // ↓ログを記録し、エラーを返す
		slog.Error("クエリーでエラーが発生", "cause", err)
		return nil, fmt.Errorf("クエリーでエラーが発生")
	}

	return struct{}{}, nil
}
```

##### log

```
{"time":"2024-06-19T12:48:11.077021+09:00","level":"ERROR","source":{"function":"main.service","file":"/Users/hotsumakawashima/git/go-stack-traces/main.go","line":34},"msg":"クエリーでエラーが発生","cause":"エラー"}
{"time":"2024-06-19T12:48:11.077618+09:00","level":"ERROR","source":{"function":"main.handler","file":"/Users/hotsumakawashima/git/go-stack-traces/main.go","line":26},"msg":"エラーが発生しました","cause":"クエリーでエラーが発生"}
```

## good

ログを記録するか、エラーを返すかのどちらかにし、両方同時に行わないようにします。
また、追加の情報が必要であれば、文脈情報を追加してエラーをラップして返すようにします。
ログの出力場所は、エントリーポイントに近いところで一度だけ行うようにします。

ここでは、スタックトレースをログに出力するために、goerr packageを使用しています。
goerr packageは、slogの構造化ログに対応しているので、ラップするだけでslogの出力にも対応できます。

##### service

```diff
func service() (interface{}, error) {
	err := query()

	if err != nil {
-       return nil, fmt.Errorf("クエリエラー")
-       slog.Error("クエリーでエラーが発生", "cause", err)
+       return nil, goerr.Wrap(err, "クエリーでエラーが発生")
	}

	return struct{}{}, nil
}
```

##### log

```
{"time":"2024-06-19T12:57:03.971038+09:00","level":"ERROR","source":{"function":"main.handler","file":"/Users/hotsumakawashima/git/go-stack-traces/main.go","line":28},"msg":"エラーが発生しました","cause":{"message":"クエリーでエラーが発生","stacktrace":["/Users/hotsumakawashima/git/go-stack-traces/main.go:36 main.service","/Users/hotsumakawashima/git/go-stack-traces/main.go:25 main.handler","/Users/hotsumakawashima/git/go-stack-traces/main.go:21 main.main","/opt/homebrew/Cellar/go/1.22.3/libexec/src/runtime/proc.go:271 runtime.main","/opt/homebrew/Cellar/go/1.22.3/libexec/src/runtime/asm_arm64.s:1222 runtime.goexit"],"cause":"エラー"}}
```

## その他の検証

### goerrのスタックトレースの出力

```json
{
  "time": "2024-06-19T12:57:03.971038+09:00",
  "level": "ERROR",
  "source": {
    "function": "main.handler",
    "file": "/Users/hotsumakawashima/git/go-stack-traces/main.go",
    "line": 28
  },
  "msg": "エラーが発生しました",
  "cause": {
    "message": "クエリーでエラーが発生",
    "stacktrace": [
      "/Users/hotsumakawashima/git/go-stack-traces/main.go:36 main.service",
      "/Users/hotsumakawashima/git/go-stack-traces/main.go:25 main.handler",
      "/Users/hotsumakawashima/git/go-stack-traces/main.go:21 main.main",
      "/opt/homebrew/Cellar/go/1.22.3/libexec/src/runtime/proc.go:271 runtime.main",
      "/opt/homebrew/Cellar/go/1.22.3/libexec/src/runtime/asm_arm64.s:1222 runtime.goexit"
    ],
    "cause": "エラー"
  }
}
```

### errの分岐

goのstandard packageでgoerrでラップされている場合も対応できることを確認。

[e9bb32b](https://github.com/hotsuma-kawashima-densan/go-stack-traces/commit/e9bb32b38403eeaf861bb3d1ccf61b415b69a920a)

```go
if errors.Is(err, SentinelError) {
    println(fmt.Sprintf("error is SentinelError: %#v", err))
}

var queryError *QueryError

if errors.As(err, &queryError) {
    println(fmt.Sprintf("unwrap error: %#v", queryError))
}
```
