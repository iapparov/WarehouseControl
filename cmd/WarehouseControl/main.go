// @title           warehouseControl API
// @version         1.0
// @description     API для управления складом.
// @BasePath        /

package main

import (
	"warehousecontrol/internal/app/history"
	"warehousecontrol/internal/app/item"
	"warehousecontrol/internal/app/user"
	"warehousecontrol/internal/auth"
	"warehousecontrol/internal/config"
	"warehousecontrol/internal/di"
	"warehousecontrol/internal/storage/postgres"
	"warehousecontrol/internal/web/handlers"

	wbzlog "github.com/wb-go/wbf/zlog"
	"go.uber.org/fx"
)

func main() {
	wbzlog.Init()
	app := fx.New(
		fx.Provide(
			config.NewAppConfig,
			postgres.NewPostgres,
			auth.NewJWTService,

			func(db *postgres.Postgres) history.HistoryStorageProvider {
				return db
			},
			history.NewHistoryService,

			func(db *postgres.Postgres) item.ItemStorageProvider {
				return db
			},
			item.NewItemService,

			func(db *postgres.Postgres) user.UserStorageProvider {
				return db
			},
			func(auth *auth.JWTService) user.JwtAuthProvider {
				return auth
			},
			user.NewUserService,

			func(app *user.UserService) handlers.UserIFace {
				return app
			},
			handlers.NewUserHandler,

			func(app *item.ItemService) handlers.ItemIFace {
				return app
			},
			handlers.NewItemHandler,

			func(app *history.HistoryService) handlers.HistoryIFace {
				return app
			},
			handlers.NewHistoryHandler,
		),
		fx.Invoke(
			di.StartHTTPServer,
			di.ClosePostgresOnStop,
		),
	)

	app.Run()
}
