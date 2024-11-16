package routes

import (
	"wallet-system/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/initTable", controllers.InitTableHandler)
	e.POST("/insertRows", controllers.InsertRowsHandler)
	e.GET("/printBalances", controllers.PrintBalancesHandler)
	e.POST("/transferFunds", controllers.TransferFundsHandler)
	e.POST("/deleteRows", controllers.DeleteRowsHandler)
}
