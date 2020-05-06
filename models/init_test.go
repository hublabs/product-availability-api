package models

import (
	"context"
	"fmt"
	"os"

	"github.com/hublabs/product-availability-api/config"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var ctx context.Context

func init() {
	os.Remove("stock.db")
	xormEngine, err := xorm.NewEngine("sqlite3", "stock.db")
	if err != nil {
		panic(fmt.Errorf("Database open error: %s \n", err))
	}
	xormEngine.ShowSQL()

	config.Init("", func(c *config.C) {
		c.Services.ColleagueAuthApi = "https://gateway-staging.srxcloud.com/srx-common/colleague-auth-api"
		c.Services.ProductApi = "https://gateway-staging.srxcloud.com/srx-brand/product-api"
	})

	Init(xormEngine)

	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine)
	behaviorlogCtx := behaviorlog.NewNopContext()
	behaviorlogCtx.AuthToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjb2xsZWFndWUiLCJjaGFubmVsQ29kZSI6IkRlZmF1bHRDb2RlIiwiY2hhbm5lbElkIjowLCJjb2xsZWFndWVOYW1lIjoi5byg5rO96Im6MTEiLCJjb2xsZWFndWVObyI6InpoYW5nLnpleWkiLCJjb2xsZWFndWVSb2xlIjoiRGVmYXVsdCIsImV4cCI6OTUzODA2ODI4MiwiaWQiOjM2NDQ3LCJpc3MiOiJjb2xsZWFndWUiLCJtYW5hZ2VkR3JvdXBJZCI6MSwibmJmIjoxNTM4MDEzOTgyLCJ0ZW5hbnRDb2RlIjoib21uaSIsInRlbmFudElkIjoxLCJ1c2VyTmFtZSI6InpoYW5nLnpleWkifQ.8sTlQefmj6J6v-EMihFMQli8rokjh4eGu1EViieou5Y"
	ctx = context.WithValue(ctx, behaviorlog.LogContextName, behaviorlogCtx)
}
