package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/hublabs/product-availability-api/factory"
)

const (
	StockTypeWarehouse = "warehouse"
	StockTypeEarlybird = "earlybird"
	StockTypeO2O       = "o2o"

	StockOperationAdd  = "ADD"
	StockOperationSet  = "SET"
	StockOperationSale = "SALE"
)

type Stock struct {
	Id        int64
	StockType string    `json:"stockType" xorm:"index"`
	SkuId     int64     `json:"skuId" xorm:"index"`
	ProductId int64     `json:"productId" xorm:"index"`
	Quantity  int64     `json:"quantity"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

type StockTotalQuantity struct {
	StockType string `json:"stockType,omitempty"`
	ProductId int64  `json:"productId,omitempty"`
	SkuId     int64  `json:"skuId,omitempty"`
	Quantity  int64  `json:"quantity,omitempty"`
}

func (Stock) GetTotalQuantities(ctx context.Context, productIds, skuIds []int64) ([]StockTotalQuantity, error) {
	if len(productIds) == 0 && len(skuIds) == 0 {
		return nil, errors.New("'productIds' or'skuIds' are required.")
	}

	q := factory.DB(ctx).Table("stock").Select("product_id, sku_id, sum(quantity) as quantity").GroupBy("product_id, sku_id")

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

	return rows, nil

}
func (Stock) getBySkuId(ctx context.Context, skuId int64) ([]Stock, error) {
	return Stock{}.GetAll(ctx, "", nil, []int64{skuId})
}

func (Stock) GetAll(ctx context.Context, stockType string, productIds, skuIds []int64) ([]Stock, error) {
	if len(productIds) == 0 && len(skuIds) == 0 {
		return nil, errors.New("'productIds' or'skuIds' are required.")
	}

	q := factory.DB(ctx)
	if stockType != "" {
		q.Where("stock_type = ?", stockType)
	}
	if len(productIds) != 0 {
		q = q.In("product_id", productIds)
	}
	if len(skuIds) != 0 {
		q = q.In("sku_id", skuIds)
	}

	var stocks []Stock
	if err := q.Find(&stocks); err != nil {
		return nil, err
	}

	return stocks, nil
}

func (Stock) Get(ctx context.Context, stockType string, skuId int64) (bool, *Stock, error) {
	stock := Stock{StockType: stockType, SkuId: skuId}
	exist, err := factory.DB(ctx).Get(&stock)
	if err != nil {
		return false, nil, err
	}
	if !exist {
		return false, nil, nil
	}

	return true, &stock, nil
}

func (Stock) AddQuantity(ctx context.Context, stockType string, productId, skuId, quantity int64, operation string) error {
	exist, stock, err := Stock{}.Get(ctx, stockType, skuId)
	if err != nil {
		return err
	}

	if !exist {
		stock = &Stock{StockType: stockType, ProductId: productId, SkuId: skuId}
		if _, err := factory.DB(ctx).Insert(stock); err != nil {
			return err
		}
	}

	var currentQuantity int64
	switch operation {
	case StockOperationAdd:
		currentQuantity = stock.Quantity + quantity
	case StockOperationSet:
		currentQuantity = quantity
	case StockOperationSale:
		currentQuantity = stock.Quantity - quantity
	}

	if err := (StockEvent{}).Create(ctx, stock.Id, quantity, operation, currentQuantity); err != nil {
		return err
	}

	switch strings.ToUpper(operation) {
	case StockOperationAdd:
		q := "update `stock` set `quantity` = `quantity` + ? where `id` = ?"
		if _, err := factory.DB(ctx).Exec(q, quantity, stock.Id); err != nil {
			return err
		}
	case StockOperationSet:
		stock.Quantity = quantity
		if _, err := factory.DB(ctx).ID(stock.Id).Cols("quantity").Update(stock); err != nil {
			return err
		}
	case StockOperationSale:
		q := "update `stock` set `quantity` = `quantity` - ? where `id` = ?"
		if _, err := factory.DB(ctx).Exec(q, quantity, stock.Id); err != nil {
			return err
		}
	}

	return nil
}
