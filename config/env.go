package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	AppsPort      string
	ApiSecret     string
	DbAutoMigrate bool
	DbName        string
	DbHost        string
	DbPort        string
	DbUser        string
	DbPassword    string
	TokenExpired  string
}

var (
	EnvFile *Env
)

func LoadEnv() {
	_ = godotenv.Load()
	env := &Env{}
	env.AppsPort = os.Getenv("APPS_PORT")
	env.ApiSecret = os.Getenv("API_SECRET")
	env.DbAutoMigrate, _ = strconv.ParseBool(os.Getenv("DB_AUTO_CREATE"))
	env.DbName = os.Getenv("DB_NAME")
	env.DbHost = os.Getenv("DB_HOST")
	env.DbPort = os.Getenv("DB_PORT")
	env.DbUser = os.Getenv("DB_USER")
	env.DbPassword = os.Getenv("DB_PASSWORD")

	EnvFile = env
}
