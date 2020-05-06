package models

import (
	"context"
	"errors"
	"time"

	"github.com/hublabs/common/auth"
	"github.com/hublabs/product-availability-api/factory"
)

type StockEvent struct {
	Id              int64     `json:"id"`
	StockId         int64     `json:"stockId" xorm:"index"`
	Quantity        int64     `json:"quantity"`
	Operation       string    `json:"operation"`
	CurrentQuantity int64     `json:"currentQuantity"`
	ColleagueId     int64     `json:"colleagueId"`
	Username        string    `json:"username"`
	CreatedAt       time.Time `json:"createdAt" xorm:"created"`

	Stock Stock `json:"stock" xorm:"-"`
}

func (StockEvent) Create(ctx context.Context, stockId int64, quantity int64, operation string, currentQuantity int64) error {
	userClaim := auth.UserClaim{}.FromCtx(ctx)

	stockHistory := StockEvent{
		StockId:         stockId,
		Quantity:        quantity,
		ColleagueId:     userClaim.ColleagueId,
		Operation:       operation,
		CurrentQuantity: currentQuantity,
		// Username:        userClaim.Username,
	}

	if _, err := factory.DB(ctx).Insert(&stockHistory); err != nil {
		return err
	}

	return nil
}
func (StockEvent) GetAll(ctx context.Context, skuId int64, skipCount, maxResultCount int) (hasMore bool, histories []StockEvent, err error) {
	var rows []struct {
		StockEvent `xorm:"extends"`
		Stock      `xorm:"extends"`
	}
	if err := factory.DB(ctx).Table("stock_event").Select("stock_event.*, stock.*").
		Join("INNER", "stock", "stock.id = stock_event.stock_id").
		Where("sku_id = ?", skuId).
		Desc("stock_event.id").
		Limit(maxResultCount+1, skipCount).
		Find(&rows); err != nil {
		return false, nil, err
	}

	if len(rows) > maxResultCount {
		hasMore = true
		rows = rows[:len(rows)-1]
	}

	for _, row := range rows {
		row.StockEvent.Stock = row.Stock
		histories = append(histories, row.StockEvent)
	}

	return hasMore, histories, nil
}
func (StockEvent) GetTotalQuantity(ctx context.Context, productIds, skuIds []int64) ([]StockTotalQuantity, error) {
	if len(productIds) == 0 && len(skuIds) == 0 {
		return nil, errors.New("'productIds' or'skuIds' are required.")
	}

	q := factory.DB(ctx).Table("stock_event").Select("stock_type, product_id, sku_id, sum(quantity) as quantity").
		Join("INNER", "stock", "stock.id = stock_event.stock_id").
		GroupBy("stock_type, product_id, sku_id")
	if len(productIds) != 0 {
		q = q.In("product_id", productIds)
	}
	if len(skuIds) != 0 {
		q = q.In("sku_id", skuIds)
	}

	var rows []StockTotalQuantity
	if err := q.Find(&rows); err != nil {
		return nil, err
	}

	if len(rows) == 1 && rows[0].Quantity == 0 {
		return nil, nil
	}

	return rows, nil
}
