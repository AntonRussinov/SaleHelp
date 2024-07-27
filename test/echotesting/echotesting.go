package echotesting

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"salehelp/pkg/extendedvalidator"

	"github.com/labstack/echo/v4"
)

type Suite struct {
	E       *echo.Echo
	Req     *http.Request
	Rec     *httptest.ResponseRecorder
	EchoCtx echo.Context
}

func New(ctx context.Context, method string, body io.Reader) Suite {
	e := echo.New()
	e.Validator = extendedvalidator.MustNew()

	req := httptest.NewRequest(method, "/", body).WithContext(ctx)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)

	return Suite{
		E:       e,
		Req:     req,
		Rec:     rec,
		EchoCtx: echoCtx,
	}
}

func (suite *Suite) SetPathParam(name, value string) {
	suite.EchoCtx.SetParamNames(name)
	suite.EchoCtx.SetParamValues(value)
}

func (suite *Suite) SetContext(ctx context.Context) {
	suite.Req = suite.EchoCtx.Request().WithContext(ctx)
	suite.EchoCtx.SetRequest(suite.Req)
}
