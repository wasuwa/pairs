package main

import (
	"pairs/pkg/auth"
	"pairs/pkg/env"
	"pairs/pkg/footprint"
	"pairs/pkg/logging"
	"pairs/pkg/selenium"

	"go.uber.org/zap"
)

func main() {
	logging.Init()
	defer logging.Sync()

	if err := env.Init(); err != nil {
		logging.Panic("環境変数が読み込めませんでした", zap.Error(err))
	}

	s, err := selenium.InitChrome()
	if err != nil {
		logging.Panic("ChromeDriverが起動できませんでした", zap.Error(err))
	}
	defer s.Stop()

	logging.Info("ブースト中の足跡付与を開始します")

	if err := auth.LoginEmail(s); err != nil {
		logging.Panic("Pairsにメールアドレスでログインできませんでした", zap.Error(err))
	}

	f, err := footprint.NewFootprint(18, 26, []string{"新潟", "富山", "石川", "埼玉", "東京", "神奈川", "千葉"}, "オンライン")
	if err != nil {
		logging.Panic("Footprintオブジェクトのバリデーションに失敗しました", zap.Error(err))
	}

	if err := f.Filtering(s); err != nil {
		logging.Panic("検索条件のフィルタリングに失敗しました", zap.Error(err))
	}

	logging.Info("ブースト中の足跡付与を終了します")
}
