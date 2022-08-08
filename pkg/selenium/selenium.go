package selenium

import "github.com/sclevine/agouti"

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
		panic(err)
	}

	page, err := driver.NewPage()
	if err != nil {
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
		panic(err)
	}
}
