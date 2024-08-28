package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"req3rdPartyServices/configs"
	"req3rdPartyServices/handler"
	"req3rdPartyServices/metrics"
	"req3rdPartyServices/repository"
	"req3rdPartyServices/service"
	"req3rdPartyServices/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logFile, err := logger.InitLogger()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	logrus.Info("logFile initialized successfully")
	defer logFile.Close()

	cfg, err := configs.InitConfig()
	if err != nil {
		logrus.WithError(err).Fatal("error initializing configs")
	}
	logrus.Info("Configs initialized successfully")

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize db")
	}
	logrus.Info("Database connected successfully")

	redisClient, err := repository.NewRedisDB(repository.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize redis")
	}
	logrus.Info("Redis connected successfully")

	go func() {
		err := metrics.StartMetricsServer(":9090")
		if err != nil {
			logrus.WithError(err).Fatal("Failed to start metrics server")
		}
		logrus.Info("Start metrics server")
	}()

	timeDuration := 10*time.Minute
	repos := repository.NewTaskRepository(db)
	services := service.NewTaskService(repos, redisClient, timeDuration)
	handlers := handler.NewTaskHandler(services)

	r := gin.Default()
	r.POST("/task", handlers.CreateTask)
	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/task/:id", handlers.GetTask)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logrus.Info("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("failed to shut down server")
	}
	logrus.Info("Server shut down successfully")
}
