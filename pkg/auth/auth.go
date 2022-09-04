package auth

import (
	"os"

	"pairs/pkg/google"
	"pairs/pkg/logging"
	"pairs/pkg/selenium"
	"time"
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
	code := message.Snippet[16:22]

	// 認証コードの入力
	fields := page.AllByClass("css-xavfkq")
	fields.At(0).Fill(code[])
	fields.At(1).Fill()
	fields.At(2).Fill()
	fields.At(3).Fill()
	fields.At(4).Fill()
	fields.At(5).Fill()
	fields.At(6).Fill()


	return nil
}
