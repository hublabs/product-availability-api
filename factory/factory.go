package factory

import (
	"context"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

var (
	_db                  *xorm.Engine
	_productDB           *xorm.Engine
	ProductDBContextName = "productDB"
)

func DB(ctx context.Context) xorm.Interface {
	v := ctx.Value(echomiddleware.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	db, ok := v.(xorm.Interface)
	if !ok {
		panic("DB is not exist")
	}
	return db
}

func ProductDB(ctx context.Context) xorm.Interface {
	v := ctx.Value(ProductDBContextName)
	if v == nil {
		panic("DB is not exist")
	}
	db, ok := v.(xorm.Interface)
	if !ok {
		panic("DB is not exist")
	}
	return db
}
