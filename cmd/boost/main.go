package main

import (
	"os"
	"pairs/pkg/auth"
	"pairs/pkg/env"
	"pairs/pkg/logging"
	"pairs/pkg/selenium"

	"go.uber.org/zap"
)

func main() {
	logging.Init()
	defer logging.Sync()
	logging.Info("ブースト中の足跡付与を開始します")

	env.Init()

	s := selenium.InitChrome()
	defer s.Stop()

	if err := auth.LoginEmail(s); err != nil {
		logging.Error("Pairsにメールアドレスでログインできませんでした", zap.Error(err))
		os.Exit(1)
	}
}
