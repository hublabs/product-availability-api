package models

import (
	"context"
	"errors"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"
)

var (
	NotFoundError = errors.New("Not found.")
)

func Init(xormEngine *xorm.Engine) {
	xormEngine.Sync(
		new(Stock),
		new(StockEvent),
	)

	ctx := context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine)

	if token, err := jwtutil.NewToken(map[string]interface{}{}); err == nil {
		behaviorlogCtx := behaviorlog.NewNopContext()
		behaviorlogCtx.AuthToken = "Bearer " + token
		ctx = context.WithValue(ctx, behaviorlog.LogContextName, behaviorlogCtx)
	}
}
