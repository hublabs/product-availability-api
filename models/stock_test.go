package models

import (
	"fmt"
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

func TestStock(t *testing.T) {
	err := Stock{}.AddQuantity(ctx, StockTypeWarehouse, 1, 1, 10, "ADD")
	test.Ok(t, err)

	err = Stock{}.AddQuantity(ctx, StockTypeWarehouse, 1, 1, 10, "ADD")
	test.Ok(t, err)

	err = Stock{}.AddQuantity(ctx, StockTypeWarehouse, 1, 1, 10, "SET")
	test.Ok(t, err)

	exist, stock, err := Stock{}.Get(ctx, StockTypeWarehouse, 1)
	test.Equals(t, exist, true)
	test.Ok(t, err)

	fmt.Println(stock.Quantity)
}
