package main

import (
	"fmt"
	"net/http"
	"pairs/pkg/logging"

	"go.uber.org/zap"
)

func main() {
	logging.Init()
	defer logging.Sync()

	logging.Info("コールバックサーバーを起動します")

	http.HandleFunc("/auth/callback/google", callbackGoogle)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logging.Error("コールバックサーバーの起動に失敗しました", zap.Error(err))
	}
}

func callbackGoogle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "このページを閉じてください")
	logging.Info("Googleの認可コード", zap.String("Code", r.FormValue("code")))
}
