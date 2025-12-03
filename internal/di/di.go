package di

import (
	"context"
	"fmt"
	"log"
	"net/http"

	wbgin "github.com/wb-go/wbf/ginext"
	"go.uber.org/fx"

	"warehousecontrol/internal/config"
	"warehousecontrol/internal/storage/postgres"
	"warehousecontrol/internal/web/handlers"
	"warehousecontrol/internal/web/routers"
)

func StartHTTPServer(lc fx.Lifecycle, userHandler *handlers.UserHandler, itemHandler *handlers.ItemHandler, historyHandler *handlers.HistoryHandler, config *config.AppConfig) {
	router := wbgin.New(config.GinConfig.Mode)

	router.Use(wbgin.Logger(), wbgin.Recovery())
	router.Use(func(c *wbgin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	routers.RegisterRoutes(router, userHandler, itemHandler, historyHandler)

	addres := fmt.Sprintf("%s:%d", config.ServerConfig.Host, config.ServerConfig.Port)
	server := &http.Server{
		Addr:    addres,
		Handler: router.Engine,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server started")
			go func() {
				if err := server.ListenAndServe(); err != nil {
					log.Printf("ListenAndServe error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Shutting down server...")
			return server.Close()
		},
	})
}

func ClosePostgresOnStop(lc fx.Lifecycle, postgres *postgres.Postgres) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing Postgres connections...")
			if err := postgres.Close(); err != nil {
				log.Printf("Failed to close Postgres: %v", err)
				return err
			}
			log.Println("Postgres closed successfully")
			return nil
		},
	})
}
