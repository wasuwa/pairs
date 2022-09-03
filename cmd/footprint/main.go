package main

import (
	"pairs/pkg/env"
	"pairs/pkg/logging"
	"pairs/pkg/selenium"
)

func main() {
	logging.Init()
	defer logging.Sync()

	env.Init()

	s := selenium.InitChrome()
	defer s.Stop()
}
