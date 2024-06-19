package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/m-mizutani/goerr"
)

func init() {
	// ログ設定
	sLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,           // ソースの行も出力
		Level:     slog.LevelInfo, // ログレベル
	}))
	slog.SetDefault(sLogger)
}

func main() {
	handler()
}

func handler() {
	_, err := service(1)

	if err != nil {
		slog.Error("1 エラーが発生しました", "cause", err)
	}

	if errors.Is(err, SentinelError) {
		println(fmt.Sprintf("error is SentinelError: %#v", err))
	}

	_, err = service(2)

	if err != nil {
		slog.Error("2 エラーが発生しました", "cause", err)
	}

	var queryError *QueryError

	if errors.As(err, &queryError) {
		println(fmt.Sprintf("unwrap error: %#v", queryError))
	}
}

type Param int

func service(param Param) (interface{}, error) {
	err := query(param)

	if err != nil {
		return nil, goerr.Wrap(err, "クエリーでエラーが発生")
	}

	return struct{}{}, nil
}

var SentinelError = errors.New("query function error")

type QueryError struct {
	message string
}

func (e *QueryError) Error() string {
	return e.message
}

func query(param Param) error {
	switch param {
	case 1:
		return SentinelError
	case 2:
		return &QueryError{message: "クエリーでエラーが発生"}
	default:
		panic("unexpected")
	}
}
