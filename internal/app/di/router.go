package di

import (
	"fmt"
	"reflect"
	"salehelp/internal/handler"
	"salehelp/internal/handler/alert"
	"salehelp/internal/handler/product"
	auth_google_handler "salehelp/internal/handler/auth_google"
	shopaisettings "salehelp/internal/handler/shop_ai_settings"
	"salehelp/internal/handler/subscription"
	"salehelp/internal/handler/telegram"
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
		router.Use(echomw.RequestLogger())
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
// 		router.Use(echomw.RequestLogger())
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
	registerAuthGoogleRoutes(group, prov)
	registerTelegramRoutes(group, handlers.TelegramBot)
}

func registerSubscriptionStatusRoute(e *echo.Group, subHandler *subscription.Handler) {
	if subHandler == nil {
		return
	}
	shopSub := e.Group("/shop/:id")
	shopSub.GET("/subscription/status", subHandler.GetSubscriptionStatusByShopID)
}

func registerSubscriptionGrantRoute(e *echo.Group, subHandler *subscription.Handler) {
	if subHandler == nil {
		return
	}
	shopSub := e.Group("/shop/:id")
	shopSub.POST("/subscription/grant", subHandler.GrantSubscriptionByShopID)
}

func registerAuthGoogleRoutes(e *echo.Group, prov *Provider) {
	if prov.Handlers.AuthGoogle == nil || prov.Services.AuthGoogle == nil {
		return
	}
	e.POST("/auth/google", prov.Handlers.AuthGoogle.LoginGoogle)
	e.POST("/auth/register", prov.Handlers.AuthGoogle.Register)
	e.POST("/auth/login", prov.Handlers.AuthGoogle.LoginEmail)
	protected := e.Group("", middleware.JWTAuth(prov.Services.AuthGoogle, auth_google_handler.SetJWTContext))
	protected.GET("/me", prov.Handlers.AuthGoogle.Me)
	protected.POST("/auth/link-telegram", prov.Handlers.AuthGoogle.LinkTelegram)
	protected.POST("/auth/confirm-telegram-link", prov.Handlers.AuthGoogle.ConfirmTelegramLink)
	protected.GET("/shops", prov.Handlers.Shop.ListShops)
	protected.POST("/shops", prov.Handlers.Shop.CreateShop)
	protected.GET("/shops/:shopId/products", prov.Handlers.Shop.ListShopProducts)
	protected.POST("/shop/:id/sync", prov.Handlers.Shop.SyncShopByID)
	registerAlertRoutes(protected, prov.Handlers.Alert)
	registerShopAISettingsRoutes(protected, prov.Handlers.ShopAISettings)
	registerSubscriptionStatusRoute(protected, prov.Handlers.Subscription)
	adminProtected := e.Group("", middleware.AdminAuth(prov.Services.AuthGoogle, auth_google_handler.SetJWTContext))
	registerSubscriptionGrantRoute(adminProtected, prov.Handlers.Subscription)
}

func registerTagRoutes(e *echo.Group, tagHandler *product.Handler) {
	e.GET("/products", tagHandler.GetProductByMerchantID,
		middleware.Pagination(),
		middleware.OrderBy(),
		middleware.QueryString(),
	)
}

func registerTelegramRoutes(e *echo.Group, telegramHandler *telegram.Handler) {
	e.POST("/telegram/webhook", telegramHandler.Webhook)
}

func registerAlertRoutes(e *echo.Group, alertHandler *alert.Handler) {
	shops := e.Group("/shops/:shopId")

	alerts := shops.Group("/alerts")
	alerts.POST("", alertHandler.CreateAlert)
	alerts.GET("", alertHandler.ListAlerts)
	alerts.GET("/:alertId", alertHandler.GetAlert)
	alerts.PUT("/:alertId", alertHandler.UpdateAlert)
	alerts.DELETE("/:alertId", alertHandler.DeleteAlert)
}

func registerShopAISettingsRoutes(e *echo.Group, aiSettingsHandler *shopaisettings.Handler) {
	if aiSettingsHandler == nil {
		return
	}
	shops := e.Group("/shops/:shopId")
	shops.GET("/ai-settings", aiSettingsHandler.GetAISettings)
	shops.PUT("/ai-settings", aiSettingsHandler.UpsertAISettings)
}
