package main

//go:generate swag init -g main.go -o ../../docs --parseDependency --parseInternal

// @title          SaleHelp API
// @version        1.0
// @description    API для сервиса SaleHelp: мониторинг магазинов Wildberries, алерты, уведомления, AI-ответы на отзывы.
// @host           localhost:8080
// @BasePath       /api
// @schemes        http https
// @Accept         json
// @Produce        json

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Точка входа для генерации Swagger-документации (swag init).
}
