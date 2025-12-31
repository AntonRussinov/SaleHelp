package di

import (
	"fmt"
	"reflect"
	"salehelp/internal/handler"
	"salehelp/internal/handler/product"
	"salehelp/internal/handler/shop"
	"salehelp/internal/middleware"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	echoswag "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (prov *Provider) InitRouter() {
	router := echo.New()
	router.Debug = prov.cfg.Environment != "production"
	router.Validator = prov.Validator
	router.HTTPErrorHandler = handler.ErrorHandler(prov.log)

	prov.registerRootMiddlewares(router)
	prov.registerSwaggerRoutes(router)

	api := router.Group("/api")
	api.GET("/version", handler.HealthCheck)
	api.Static("/assets", "assets")

	v1 := api.Group("/v1")
	prov.registerV1Routes(v1)

	prov.Router = router
}

func (prov *Provider) registerRootMiddlewares(router *echo.Echo) {
	router.Use(
		echomw.Recover(),
		echomw.CORSWithConfig(
			echomw.CORSConfig{
				AllowOrigins: prov.cfg.AllowedOrigins,
				AllowMethods: echomw.DefaultCORSConfig.AllowMethods,
				AllowHeaders: []string{
					echo.HeaderAuthorization, echo.HeaderAccept,
					echo.HeaderOrigin, echo.HeaderContentType,
					"Accesstoken",
				},
				AllowCredentials: true,
			},
		),
		echomw.RequestID(),
	)
}

func (prov *Provider) registerSwaggerRoutes(router *echo.Echo) {
	if prov.cfg.Environment != "production" {
		router.GET("/swagger/*", echoswag.WrapHandler)
		router.Debug = true
		router.Use(echomw.Logger())
	}
}

func RegisterValidator(echo *echo.Echo) {
	validate := validator.New()

	_ = validate.RegisterValidation("is_valid_enum", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field()

		fieldType := fieldValue.Type()

		method, ok := fieldType.MethodByName("IsValid")
		if !ok {
			return false
		}

		output := method.Func.Call([]reflect.Value{fieldValue})
		if len(output) != 1 {
			return false
		}

		value := output[0].Interface()

		boolValue, ok := value.(bool)
		if !ok {
			return false
		}

		return boolValue
	})

	echo.Validator = &CustomValidator{Validator: validate}
}

// func NewRouter(cfg *Config, log *slog.Logger, handlers *Handler) *echo.Echo {
// 	router := echo.New()
// 	router.Debug = cfg.Environment != "production"
// 	RegisterValidator(router)
// 	router.HTTPErrorHandler = handler.ErrorHandler(log)

// 	router.Use(
// 		echomw.Recover(),
// 		echomw.CORSWithConfig(
// 			echomw.CORSConfig{
// 				AllowOrigins: cfg.AllowedOrigins,
// 				AllowMethods: echomw.DefaultCORSConfig.AllowMethods,
// 				AllowHeaders: []string{
// 					echo.HeaderAuthorization, echo.HeaderAccept,
// 					echo.HeaderOrigin, echo.HeaderContentType,
// 				},
// 				AllowCredentials: true,
// 			},
// 		),
// 		echomw.RequestID(),
// 	)

// 	if cfg.Environment != "production" {
// 		router.GET("/swagger/*", echoswag.WrapHandler)
// 		router.Debug = true
// 		router.Use(echomw.Logger())
// 	}

// 	apiGroup := router.Group("/api")
// 	//registerPublicRootRoutes(apiGroup)

// 	v1Group := apiGroup.Group("/v1")
// 	registerV1Routes(v1Group, handlers)

// 	return router
// }

func (cv *CustomValidator) Validate(i any) error {
	const op = "app.Validate"

	err := cv.Validator.Struct(i)
	if err != nil {
		return fmt.Errorf("validation failed %s: %w", op, err)
	}

	return nil
}

func (prov *Provider) registerV1Routes(group *echo.Group) {
	handlers := prov.Handlers
	group.Static("/assets", "assets")

	//protected := group.Group("", handlers.Auth.HeaderMiddleware())

	registerTagRoutes(group, handlers.Product)
	registerUserShopRoutes(group, handlers.Shop)
}

func registerUserShopRoutes(e *echo.Group, shopHandler *shop.Handler) {
	e.POST("/shop/:id/sync", shopHandler.SyncShopByID)
}

func registerTagRoutes(e *echo.Group, tagHandler *product.Handler) {
	e.GET("/products", tagHandler.GetProductByMerchantID,
		middleware.Pagination(),
		middleware.OrderBy(),
		middleware.QueryString(),
	)
}
