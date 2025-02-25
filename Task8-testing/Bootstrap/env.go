package bootstrap

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv                string `mapstructure:"APP_ENV"`
	ServerAddress         string `mapstructure:"SERVER_ADDRESS"`
	ContextTimeout        int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost                string `mapstructure:"DB_HOST"`
	DBPort                string `mapstructure:"DB_PORT"`
	DBUser                string `mapstructure:"DB_USER"`
	DBPass                string `mapstructure:"DB_PASS"`
	DBName                string `mapstructure:"DB_NAME"`
	TestDBName            string `mapstructure:"DB_NAME"`
	AccessTokenExpiryHour int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
	AccessTokenSecret     string `mapstructure:"ACCESS_TOKEN"`
}

func NewEnv(depth int) *Env {
	env := Env{}
	var path []string
	for i := 1; i < depth; i++ {
		path = append(path, "..")
	}

	projectRoot, err := filepath.Abs(filepath.Join(path...)) // Move up two directories
	if err != nil {
		log.Fatalf("Error getting project root path: %v", err)
		return nil
	}

	log.Printf("Project root path: %s", projectRoot)

	// Set the path to the .env file
	viper.SetConfigFile(filepath.Join(projectRoot, ".env"))

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Can't find the file .env: %v", err)
		return nil
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatalf("Environment can't be loaded: %v", err)
		return nil
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
