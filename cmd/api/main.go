package main

import (
	"fmt"
	stdlog "log"
	"os"
	"ps-beli-mang/cmd/api/server"
	"ps-beli-mang/configs"
	imagehandler "ps-beli-mang/internal/image/handler"
	imageservice "ps-beli-mang/internal/image/service"
	merchanthandler "ps-beli-mang/internal/merchant/handler"
	merchantrepository "ps-beli-mang/internal/merchant/repository"
	merchantservice "ps-beli-mang/internal/merchant/service"
	purchasehandler "ps-beli-mang/internal/purchase/handler"
	purchaserepository "ps-beli-mang/internal/purchase/repository"
	purchaseservice "ps-beli-mang/internal/purchase/service"
	userhandler "ps-beli-mang/internal/user/handler"
	userrepository "ps-beli-mang/internal/user/repository"
	userservice "ps-beli-mang/internal/user/service"

	bhandler "ps-beli-mang/pkg/base/handler"
	"ps-beli-mang/pkg/logger"
	psqlqgen "ps-beli-mang/pkg/psqlqgen"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var port int

var httpCmd = &cobra.Command{
	Use:   "http [OPTIONS]",
	Short: "Run HTTP API",
	Long:  "Run HTTP API for SCM",
	RunE:  runHttpCommand,
}

var (
	params          map[string]string
	baseHandler     *bhandler.BaseHTTPHandler
	userHandler     *userhandler.UserHandler
	imageHandler    *imagehandler.ImageHandler
	merchantHandler *merchanthandler.MerchantHandler
	purchaseHandler *purchasehandler.PurchaseHandler
	cfg             *configs.MainConfig
)

func init() {
	httpCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the HTTP server")
}

func main() {
	if err := httpCmd.Execute(); err != nil {
		slog.Error(fmt.Sprintf("Error on command execution: %s", err.Error()))
		os.Exit(1)
	}
}

func logLevel() slog.Level {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func initLogger() {
	{

		log, err := logger.SlogOption{
			Resource: map[string]string{
				"service.name":        "halo suster app",
				"service.ns":          "halo suster",
				"service.instance_id": "random-uuid",
				"service.version":     "v.0",
				"service.env":         "staging",
			},
			ContextExtractor:   nil,
			AttributeFormatter: nil,
			Writer:             os.Stdout,
			Leveler:            logLevel(),
		}.NewSlog()
		if err != nil {
			err = fmt.Errorf("prepare logger error: %w", err)
			stdlog.Fatal(err) // if logger cannot be prepared (commonly due to option value error), use std logger.
			return
		}

		// Set logger as global logger.
		slog.SetDefault(log)
	}
}

func runHttpCommand(cmd *cobra.Command, args []string) error {
	initLogger()
	initInfra()

	httpServer := server.NewServer(
		baseHandler,
		userHandler,
		imageHandler,
		merchantHandler,
		purchaseHandler,
		port,
	)

	return httpServer.Run()
}

func dbInitConnection() *sqlx.DB {
	return psqlqgen.Init(cfg)
}

func initInfra() {
	cfg = configs.Init()
	db := dbInitConnection()

	userRepository := userrepository.NewUserRepositoryImpl(db)
	userService := userservice.NewUserServiceImpl(userRepository)
	userHandler = userhandler.NewUserHandler(userService)
	imageService := imageservice.NewImageService(cfg)
	imageHandler = imagehandler.NewImageHandler(imageService, userService)

	merchantRepository := merchantrepository.NewMerchantRepositoryImpl(db)
	orderRepository := purchaserepository.NewOrderRepositoryImpl(db)

	merchantService := merchantservice.NewMerchantServiceImpl(merchantRepository, orderRepository)
	merchantHandler = merchanthandler.NewMerchantHandler(merchantService, userService)

	orderService := purchaseservice.NewOrderServiceImpl(orderRepository)
	purchaseHandler = purchasehandler.NewPurchaseHandler(orderService, userService)

}
