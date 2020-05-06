package models

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hublabs/product-availability-api/config"
)

type Product struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	TitleImage  string    `json:"titleImage"`
	ListPrice   float64   `json:"listPrice"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Identifiers []struct {
		Uid    string `json:"uid"`
		Source string `json:"source"`
	} `json:"identifiers"`
	Brand struct {
		Id   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"brand"`
}
type Sku struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Identifiers []struct {
		Uid    string `json:"uid"`
		Source string `json:"source"`
	} `json:"identifiers"`
	Options []struct {
		Id    int64  `json:"id"`
		SkuId int64  `json:"skuId"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"options"`
	Product Product `json:"product"`
	Plus    []struct {
		SalePrice float64 `json:"salePrice"`
		Code      string  `json:"code"`
	} `json:"plus"`
}

func (Product) GetById(ctx context.Context, producId int64) (*Product, error) {
	var resp struct {
		Result  Product `json:"result"`
		Success bool    `json:"success"`
		Error   struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}

	if err := RetryWebApi(ctx, &resp, http.MethodGet, config.Config().Services.ProductApi+"/v1/products/ids="+strconv.FormatInt(producId, 10), nil); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}

func (Sku) GetById(ctx context.Context, skuId int64) (*Sku, error) {
	var resp struct {
		Result  Sku  `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}

	if err := RetryWebApi(ctx, &resp, http.MethodGet, config.Config().Services.ProductApi+"/v1/skus/"+strconv.FormatInt(skuId, 10), nil); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result, nil
}
