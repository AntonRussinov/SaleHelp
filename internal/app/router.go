package app

import (
	"fmt"
	"log/slog"
	"reflect"
	"salehelp/internal/handler"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	echoswag "github.com/swaggo/echo-swagger"
)

type CustomValidator struct {
	Validator *validator.Validate
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

func NewRouter(cfg *Config, log *slog.Logger, handlers *Handler) *echo.Echo {
	router := echo.New()
	router.Debug = cfg.Environment != "production"
	RegisterValidator(router)
	router.HTTPErrorHandler = handler.ErrorHandler(log)

	router.Use(
		echomw.Recover(),
		echomw.CORSWithConfig(
			echomw.CORSConfig{
				AllowOrigins: cfg.AllowedOrigins,
				AllowMethods: echomw.DefaultCORSConfig.AllowMethods,
				AllowHeaders: []string{
					echo.HeaderAuthorization, echo.HeaderAccept,
					echo.HeaderOrigin, echo.HeaderContentType,
				},
				AllowCredentials: true,
			},
		),
		echomw.RequestID(),
	)

	if cfg.Environment != "production" {
		router.GET("/swagger/*", echoswag.WrapHandler)
		router.Debug = true
		router.Use(echomw.Logger())
	}

	apiGroup := router.Group("/api")
	registerPublicRootRoutes(apiGroup)

	v1Group := apiGroup.Group("/v1")
	registerV1Routes(v1Group, handlers)

	return router
}

func (cv *CustomValidator) Validate(i any) error {
	const op = "app.Validate"

	err := cv.Validator.Struct(i)
	if err != nil {
		return fmt.Errorf("validation failed %s: %w", op, err)
	}

	return nil
}

func registerPublicRootRoutes(e *echo.Group) {
	e.GET("/version", handler.HealthCheck)
	e.Static("/assets", "assets")
}

func registerV1Routes(group *echo.Group, _ *Handler) {
	group.Static("/assets", "assets")

	//registerTagRoutes(protected, handlers.Tag)
}

// func registerTagRoutes(e *echo.Group, tagHandler *tag.Handler) {
// 	e.GET("/tags", tagHandler.GeTags,
// 		middleware.Pagination(),
// 		middleware.OrderBy(),
// 		middleware.QueryString(),
// 	)
// }
