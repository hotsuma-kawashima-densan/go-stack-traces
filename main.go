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
	return nil, fmt.Errorf("エラー")
}
