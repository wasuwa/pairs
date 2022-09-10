package env

import "github.com/joho/godotenv"

// Init 初期化
func Init() error {
	if err := godotenv.Load("config/.env"); err != nil {
		return err
	}
	return nil
}
