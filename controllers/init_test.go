package controllers

import (
	"context"
	"net/http"
	"runtime"

	"github.com/hublabs/common/auth"
	"github.com/hublabs/product-availability-api/models"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	runtime.GOMAXPROCS(1)
	xormEngine, err := xorm.NewEngine("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	models.Init(xormEngine)
	xormEngine.ShowSQL()

	echoApp = echo.New()
	db := echomiddleware.ContextDB("test", xormEngine, echomiddleware.KafkaConfig{})

	behaviorlogger := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				behaviorlog.LogContextName, behaviorlog.New("test", req),
			)))
			return next(c)
		}
	}
	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return behaviorlogger(auth.UserClaimMiddleware()(db(handlerFunc)))(c)
	}
}
func setHeader(r *http.Request) {
	token, _ := jwtutil.NewToken(map[string]interface{}{"aud": "colleague"})
	r.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}
