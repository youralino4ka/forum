package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/youralino4ka/forum/internal/user"
	"github.com/youralino4ka/forum/internal/user/repository"
	"github.com/youralino4ka/forum/pkg/config"
	"github.com/youralino4ka/forum/pkg/logger"
	"github.com/youralino4ka/forum/proto/user"
)

func main() {
	// Инициализация логгера
	logger.Init(zerolog.InfoLevel)
	log := log.With().Str("service", "user-service").Logger()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Подключение к базе данных
	db, err := sqlx.Connect("postgres", cfg.DB.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Применение миграций
	if err := applyMigrations(cfg.DB.DSN, "migrations/user-service"); err != nil {
		log.Fatal().Err(err).Msg("Failed to apply migrations")
	}

	// Инициализация репозитория и сервиса
	userRepo := repository.NewPostgresRepository(db, log)
	userService := user.NewService(userRepo, log)

	// Создание GRPC сервера
	grpcServer := grpc.NewServer()
	userServer := NewUserServer(userService, log)
	userpb.RegisterUserServiceServer(grpcServer, userServer)

	// Запуск GRPC сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	go func() {
		log.Info().Msgf("Starting user service gRPC server on port %d", cfg.GRPC.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC server")
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down user service...")
	grpcServer.GracefulStop()
	log.Info().Msg("User service stopped")
}

func applyMigrations(dsn, migrationPath string) error {
	// Реализация применения миграций с помощью golang-migrate
	// В реальном проекте используйте библиотеку github.com/golang-migrate/migrate/v4
	return nil
}
