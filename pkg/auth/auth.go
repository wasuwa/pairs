package auth

import (
	"os"

	"pairs/pkg/google"
	"pairs/pkg/logging"
	"pairs/pkg/selenium"
	"time"

	"go.uber.org/zap"
)

const (
	emailLoginURL = "https://pairs.lv/login/email"
)

// LoginEmail ペアーズにログインする
func LoginEmail(selenium *selenium.Selenium) error {
	logging.Info("Pairsにログインします")

	page := selenium.Page
	if err := page.Navigate(emailLoginURL); err != nil {
		return err
	}

	// メールアドレスの入力
	field := page.FindByName("email")
	field.Fill(os.Getenv("PAIRS_EMAIL"))
	page.FindByClass("css-hjwi6k").Click()
	time.Sleep(time.Second * 2)

	page.FindByClass("css-i10lev").Click()
	time.Sleep(time.Second * 2)

	// 認証コードの取得
	gmail, err := google.NewGmail()
	if err != nil {
		return err
	}
	message, err := gmail.FetchEmail("https://pairs.lv/app/mail_auth")
	if err != nil {
		return err
	}

	// 認証コードの入力
	code := message.Snippet[16:22]
	logging.Info("ペアーズの認証コードをブラウザで入力してください（10秒以内）", zap.String("AuthCode", code))
	time.Sleep(time.Second * 10)

	return nil
}
