package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/hublabs/common/api"
	productmodule "github.com/hublabs/product-api/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func TestStock(t *testing.T) {
	p := productmodule.Product{
		Id: 100,
		Skus: []productmodule.Sku{
			{Id: 200},
			{Id: 201},
		},
	}
	t.Run("Set Stock", func(t *testing.T) {
		param := []map[string]interface{}{
			{
				"stockType": "warehouse",
				"productId": p.Id,
				"skuId":     p.Skus[0].Id,
				"quantity":  5,
				"operation": "SET",
			},
			{
				"stockType": "o2o",
				"productId": p.Id,
				"skuId":     p.Skus[0].Id,
				"quantity":  5,
				"operation": "SET",
			},
			{
				"stockType": "warehouse",
				"productId": p.Id,
				"skuId":     p.Skus[1].Id,
				"quantity":  5,
				"operation": "SET",
			},
			{
				"stockType": "o2o",
				"productId": p.Id,
				"skuId":     p.Skus[1].Id,
				"quantity":  5,
				"operation": "SET",
			},
		}

		body, _ := json.Marshal(param)

		req := httptest.NewRequest(echo.POST, "/v1/stocks", bytes.NewReader(body))
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		test.Ok(t, handleWithFilter(StockController{}.AddQuantity, c))
		test.Equals(t, http.StatusOK, rec.Code)
	})

	t.Run("Check Stock By ProductId", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/stocks/:productId", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetParamNames("productId")
		c.SetParamValues(strconv.FormatInt(p.Id, 10))
		test.Ok(t, handleWithFilter(StockController{}.GetAvailableByProductId, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Success bool `json:"success"`
			Result  []struct {
				SkuId             int64
				ProductId         int64
				AvailableQuantity int
			} `json:"result"`
			Error api.Error `json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Success, true)
		test.Equals(t, v.Error.Message, "")
		for _, q := range v.Result {
			test.Equals(t, q.ProductId, p.Id)
			test.Equals(t, q.AvailableQuantity, 10)
		}
	})

	t.Run("Check Stock By SkuId", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/stocks/:productId/:skuId", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetParamNames("productId", "skuId")
		c.SetParamValues(strconv.FormatInt(p.Id, 10), strconv.FormatInt(p.Skus[0].Id, 10))
		test.Ok(t, handleWithFilter(StockController{}.GetAvailableByProductIdAndSkuId, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Success bool `json:"success"`
			Result  struct {
				SkuId             int64
				ProductId         int64
				AvailableQuantity int
			} `json:"result"`
			Error api.Error `json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Success, true)
		test.Equals(t, v.Error.Message, "")
		test.Equals(t, v.Result.AvailableQuantity, 10)
	})

	t.Run("Check Stock By ProductId", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/stocks/:productId", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetParamNames("productId")
		c.SetParamValues(strconv.FormatInt(p.Id, 10))
		test.Ok(t, handleWithFilter(StockController{}.GetAvailableByProductId, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Success bool `json:"success"`
			Result  []struct {
				SkuId             int64
				ProductId         int64
				AvailableQuantity int
			} `json:"result"`
			Error api.Error `json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Success, true)
		test.Equals(t, v.Error.Message, "")
		for _, q := range v.Result {
			test.Equals(t, q.ProductId, p.Id)
			test.Equals(t, q.AvailableQuantity, 10)
		}
	})

	t.Run("Check Stock By SkuId", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/stocks/:productId/:skuId", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetParamNames("productId", "skuId")
		c.SetParamValues(strconv.FormatInt(p.Id, 10), strconv.FormatInt(p.Skus[0].Id, 10))
		test.Ok(t, handleWithFilter(StockController{}.GetAvailableByProductIdAndSkuId, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Success bool `json:"success"`
			Result  struct {
				SkuId             int64
				ProductId         int64
				AvailableQuantity int
			} `json:"result"`
			Error api.Error `json:"error"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Success, true)
		test.Equals(t, v.Error.Message, "")
		test.Equals(t, v.Result.AvailableQuantity, 10)
	})
}
