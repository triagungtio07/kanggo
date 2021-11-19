package main

import (
	"fmt"
	"kanggo/config"
	userHandler "kanggo/pkg/handler/user"
	userStorage "kanggo/pkg/storage/user"
	userUsecase "kanggo/pkg/usecase/user"

	productHandler "kanggo/pkg/handler/product"
	productStorage "kanggo/pkg/storage/product"
	productUsecase "kanggo/pkg/usecase/product"

	orderHandler "kanggo/pkg/handler/order"
	orderStorage "kanggo/pkg/storage/order"
	orderUsecase "kanggo/pkg/usecase/order"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": 200, "message": "All good"})
	})

	config.ConnectDb()

	//storage
	userStorage := userStorage.NewUserStorage(config.Native, config.Gorm)
	productStorage := productStorage.NewProductStorage(config.Native, config.Gorm)
	orderStorage := orderStorage.NewOrderStorage(config.Native, config.Gorm)

	//usecase
	userUsecase := userUsecase.NewUserUsecase(userStorage)
	productUsecase := productUsecase.NewProductUsecase(productStorage)
	orderUsecase := orderUsecase.NewOrderUsecase(orderStorage, productStorage)

	//handler
	userHandler := userHandler.NewUserhandler(userUsecase)
	productHandler := productHandler.NewProductHandler(productUsecase)
	orderHandler := orderHandler.NewOrderHandler(orderUsecase)

	//router
	userHandler.Route(engine)
	productHandler.Route(engine)
	orderHandler.Route(engine)

	fmt.Println("Running on port : 8080")
	engine.Run(config.EnvFile.AppsPort)

}
