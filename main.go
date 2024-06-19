package main

import (
	"fmt"
	"log/slog"
	"os"
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
	_, err := service()

	if err != nil {
		slog.Error("エラーが発生しました", "cause", err)
	}
}

func service() (interface{}, error) {
	err := query()

	if err != nil {
		slog.Error("クエリーでエラーが発生", "cause", err)
		return nil, fmt.Errorf("クエリーでエラーが発生")
	}

	return struct{}{}, nil
}

func query() error {
	return fmt.Errorf("エラー")
}
