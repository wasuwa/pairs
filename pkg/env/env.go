package env

import "github.com/joho/godotenv"

// Init 初期化
func Init() {
	if err := godotenv.Load("config/.env"); err != nil {
		panic(err)
	}
}
