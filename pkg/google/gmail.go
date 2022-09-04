package google

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pairs/pkg/logging"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	errEmptyMessages  = errors.New("メッセージが空です")
	errNoMatchKeyword = errors.New("キーワードに一致するメールが見つかりませんでした")
)

// Gmail Google Email
type Gmail struct {
	Service *gmail.Service
}

// NewGmail 初期化
func NewGmail() (*Gmail, error) {
	logging.Info("Gmailに接続します")

	ctx := context.Background()
	config := oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/callback/google",
		Scopes:       []string{"https://mail.google.com/"},
	}

	url := config.AuthCodeURL("state")
	logging.Info("Gmailの認証URLです（ペアーズのメールアドレスでログインしてください）", zap.String("URL", url))

	var code string
	fmt.Print("URLに含まれている「code」を入力してください: ")
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	service, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}

	return &Gmail{Service: service}, nil
}

// FetchEmail メールを取得する
func (g *Gmail) FetchEmail(keyword string) (*gmail.Message, error) {
	messages := g.Service.Users.Messages
	res, err := messages.List("me").Do()
	if err != nil {
		return nil, err
	}
	if len(res.Messages) <= 0 {
		return nil, errEmptyMessages
	}

	for _, v := range res.Messages {
		m, err := messages.Get("me", v.Id).Do()
		if err != nil {
			return nil, err
		}
		// 本文にキーワードが含まれていた場合
		if strings.Contains(m.Snippet, keyword) {
			return m, nil
		}
	}
	// キーワードとマッチしない場合はエラー
	return nil, errNoMatchKeyword
}
