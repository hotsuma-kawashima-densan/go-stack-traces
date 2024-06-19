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

	err := fmt.Errorf("エラー")

	slog.Error("ログ出力", "cause", err)
}
