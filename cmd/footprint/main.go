package main

import (
	"pairs/pkg/logging"
	"pairs/pkg/selenium"
)

func main() {
	logging.Init()
	defer logging.Sync()

	s := selenium.InitChrome()
	defer s.Stop()
}
