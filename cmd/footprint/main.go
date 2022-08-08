package main

import "pairs/pkg/selenium"

func main() {
	s := selenium.InitChrome()
	defer s.Stop()
}
