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
func InitChrome() (*Selenium, error) {
	logging.Info("ChromeDriverを起動します")

	opts := agouti.ChromeOptions(
		"args", []string{
			"--disable-gpu",
		},
	)
	driver := agouti.ChromeDriver(opts)
	if err := driver.Start(); err != nil {
		return nil, err
	}

	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}

	s := &Selenium{
		Driver: driver,
		Page:   page,
	}
	return s, nil
}

// Stop Web driver の停止
func (s *Selenium) Stop() {
	if err := s.Driver.Stop(); err != nil {
		logging.Panic("ChromeDriverが停止できませんでした", zap.Error(err))
	}
}
