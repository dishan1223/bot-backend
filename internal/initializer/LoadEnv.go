package initializer

import "github.com/joho/godotenv"

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Internal SERVER Error :: Failed to Load Environments.")
	}
}
