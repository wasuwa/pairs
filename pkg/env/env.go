package env

import (
	"pairs/pkg/logging"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Init 初期化
func Init() {
	if err := godotenv.Load("config/.env"); err != nil {
		logging.Panic("環境変数が読み込めませんでした", zap.Error(err))
		panic(err)
	}
}
