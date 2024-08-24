package main

import (
	"req3rdPartyServices/configs"
	"req3rdPartyServices/handler"
	"req3rdPartyServices/repository"
	"req3rdPartyServices/service"
	"req3rdPartyServices/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logFile, err := logger.InitLogger()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	logrus.Info("logFile initialized successfully")
	defer logFile.Close()

	if err := configs.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing configs")
	}
	logrus.Info("Configs initialized successfully")

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize db")
	}
	logrus.Info("Database connected successfully")

	redisClient, err := repository.NewRedisDB(repository.RedisConfig{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.DB"),
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize redis")
	}
	logrus.Info("Redis connected successfully")

	repos := repository.NewTaskRepository(db)
	services := service.NewTaskService(repos, redisClient, 10*time.Minute)
	handlers := handler.NewTaskHandler(services)

	r := gin.Default()
	r.POST("/task", handlers.CreateTask)
	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/task/:id", handlers.GetTask)

	if err := r.Run(":8080"); err != nil {
		logrus.Fatal("failed to start server: ", err.Error())
	}
	logrus.Info("The server has been started successfully")
}
