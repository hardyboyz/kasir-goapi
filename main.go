package main

import (
	"log"
	database "myshop-api/config"
	"myshop-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// Load configuration
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("Unable to decode into struct, %v", err)
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Setup Gin router
	r := gin.Default()
	routes.SetupRoutes(r, db)

	log.Printf("Server starting on port %s", config.Port)
	r.Run(config.Port)
}
