package selenium

import (
	"pairs/pkg/logging"

	"github.com/sclevine/agouti"
	"go.uber.org/zap"
)

// Selenium selenium
type Selenium struct {
	Page   *agouti.Page
	Driver *agouti.WebDriver
}

// InitChrome Chrome Web driver のセットアップ
func InitChrome() *Selenium {
	opts := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--disable-gpu",
		},
	)
	driver := agouti.ChromeDriver(opts)
	if err := driver.Start(); err != nil {
		logging.Panic("ChromeDriverが起動できません", zap.Error(err))
		panic(err)
	}

	page, err := driver.NewPage()
	if err != nil {
		logging.Panic("Pageの初期化に失敗しました", zap.Error(err))
		panic(err)
	}

	return &Selenium{
		Driver: driver,
		Page:   page,
	}
}

// Stop Web driver の停止
func (s *Selenium) Stop() {
	if err := s.Driver.Stop(); err != nil {
		logging.Panic("ChromeDriverが停止できません", zap.Error(err))
		panic(err)
	}
}
