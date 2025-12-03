package routers

import (
	_ "warehousecontrol/docs"
	"warehousecontrol/internal/domain/user"
	"warehousecontrol/internal/web/handlers"

	httpSwagger "github.com/swaggo/http-swagger"
	wbgin "github.com/wb-go/wbf/ginext"
)

func RegisterRoutes(engine *wbgin.Engine, userHandler *handlers.UserHandler, itemHandler *handlers.ItemHandler, historyHandler *handlers.HistoryHandler) {
	api := engine.Group("/api")
	api.GET("/swagger/*any", func(c *wbgin.Context) {
		httpSwagger.WrapHandler(c.Writer, c.Request)
	})

	// публичные маршруты авторизации
	auth := api.Group("/auth")
	auth.POST("/register", userHandler.RegisterUser)
	auth.POST("/login", userHandler.LoginUser)
	auth.POST("/refresh", userHandler.RefreshToken)

	// защищённая группа предметов
	items := api.Group("/items", AuthMiddleware(userHandler.Service))
	items.POST("", RequireRoles(user.Admin), itemHandler.CreateItem)
	items.GET("", RequireRoles(user.Admin, user.Manager, user.Viewer), itemHandler.GetItems)
	items.GET("/:id", RequireRoles(user.Admin, user.Manager, user.Viewer), itemHandler.GetItem)
	items.PUT("/:id", RequireRoles(user.Admin, user.Manager), itemHandler.PutItem)
	items.DELETE("/:id", RequireRoles(user.Admin), itemHandler.DeleteItem)

	// просмотр истории только для админа
	history := api.Group("/history", AuthMiddleware(userHandler.Service), RequireRoles(user.Admin))
	history.GET("", historyHandler.GetItems)
	history.GET("/csv", historyHandler.GetItemsCSV)
}
