package controllers

import (
	"net/http"
	"strconv"

	"github.com/hublabs/common/api"
	"github.com/hublabs/product-availability-api/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type StockController struct{}

func (c StockController) Init(g echoswagger.ApiGroup) {
	g.GET("/histories", c.GetHistories).
		AddParamQuery(1, "skuId", "sku id", true).
		AddParamQuery(1, "skipCount", "skip count", false).
		AddParamQuery(1, "maxResultCount", "max result count", false).
		SetSecurity("Authorization")
	g.GET("", c.GetAll).
		AddParamQuery("", "skuIds", "sku id list", false).
		AddParamQuery("", "productIds", "product id list", false).
		SetSecurity("Authorization")

	g.POST("", c.AddQuantity).
		AddParamBody([]AddQuantityParam{}, "AddQuantityParam", "AddQuantityParam", true).
		SetSecurity("Authorization")

	// available quantity = stock quantity + earlybird quantity
	g.GET("/:productId", c.GetAvailableByProductId).
		AddParamPath(1, "productId", "product id").
		SetSecurity("Authorization").
		SetSecurity("Authorization")

	// available quantity = stock quantity + earlybird quantity
	g.GET("/:productId/:skuId", c.GetAvailableByProductIdAndSkuId).
		AddParamPath(1, "productId", "product id").
		AddParamPath(1, "skuId", "sku id").
		SetSecurity("Authorization")
}
func (StockController) GetHistories(c echo.Context) error {
	maxResultCount, _ := strconv.Atoi(c.QueryParam("maxResultCount"))
	skipCount, _ := strconv.Atoi(c.QueryParam("skipCount"))
	skuId, _ := strconv.ParseInt(c.QueryParam("skuId"), 10, 64)

	if skuId == 0 {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{Message: "Invalid skuId"},
		})
	}
	if maxResultCount == 0 {
		maxResultCount = defaultMaxResultCount
	}

	hasMore, histories, err := models.StockEvent{}.GetAll(c.Request().Context(), skuId, skipCount, maxResultCount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Result{
			Success: false,
			Error: api.Error{
				Message: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, api.Result{
		Success: true,
		Result: map[string]interface{}{
			"hasMore": hasMore,
			"items":   histories,
		},
	})

}
func (StockController) GetAll(c echo.Context) error {
	skuIds := convertStrToArrInt64(c.QueryParam("skuIds"))
	productIds := convertStrToArrInt64(c.QueryParam("productIds"))
	if len(productIds) == 0 && len(skuIds) == 0 {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{
				Message: "'productIds' or'skuIds' are required.",
			},
		})
	}

	stocks, err := models.Stock{}.GetAll(c.Request().Context(), c.QueryParam("stockType"), productIds, skuIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Result{
			Error: api.Error{
				Message: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, api.Result{
		Success: true,
		Result:  stocks,
	})
}

// Available Quantity = Stock Quantity + Earlybird Quantity
func (StockController) GetAvailableByProductId(c echo.Context) error {
	productId, _ := strconv.ParseInt(c.Param("productId"), 10, 64)
	if productId == 0 {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{
				Message: "Invalid productId",
			},
		})
	}

	totalQuantities, err := models.Stock{}.GetTotalQuantities(c.Request().Context(), []int64{productId}, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Result{
			Error: api.Error{
				Message: err.Error(),
			},
		})
	}

	result := []map[string]int64{}
	for _, s := range totalQuantities {
		result = append(result, map[string]int64{
			"skuId":             s.SkuId,
			"productId":         s.ProductId,
			"availableQuantity": s.Quantity,
		})
	}

	return c.JSON(http.StatusOK, api.Result{
		Success: true,
		Result:  result,
	})
}

// Available Quantity = Stock Quantity + Earlybird Quantity
func (StockController) GetAvailableByProductIdAndSkuId(c echo.Context) error {
	productId, _ := strconv.ParseInt(c.Param("productId"), 10, 64)
	if productId == 0 {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{
				Message: "Invalid productId",
			},
		})
	}

	skuId, _ := strconv.ParseInt(c.Param("skuId"), 10, 64)
	if skuId == 0 {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{
				Message: "Invalid skuId",
			},
		})
	}

	stocks, err := models.Stock{}.GetTotalQuantities(c.Request().Context(), nil, []int64{skuId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, api.Result{
			Error: api.Error{
				Message: err.Error(),
			},
		})
	}

	stockTotalQuantity := models.StockTotalQuantity{
		SkuId:     skuId,
		ProductId: productId,
	}
	if len(stocks) != 0 {
		stockTotalQuantity.Quantity = stocks[0].Quantity
	}

	return c.JSON(http.StatusOK, api.Result{
		Success: true,
		Result: map[string]int64{
			"skuId":             stockTotalQuantity.SkuId,
			"productId":         stockTotalQuantity.ProductId,
			"availableQuantity": stockTotalQuantity.Quantity,
		},
	})
}

type AddQuantityParam struct {
	StockType string `json:"stockType"`
	ProductId int64  `json:"productId"`
	SkuId     int64  `json:"skuId"`
	Quantity  int64  `json:"quantity"`
	Operation string `json:"operation"`
}

func (StockController) AddQuantity(c echo.Context) error {
	var params []AddQuantityParam
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, api.Result{
			Error: api.Error{
				Message: err.Error(),
			},
		})
	}

	for _, v := range params {
		err := models.Stock{}.AddQuantity(c.Request().Context(), v.StockType, v.ProductId, v.SkuId, v.Quantity, v.Operation)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, api.Result{
				Error: api.Error{
					Message: err.Error(),
				},
			})
		}
	}

	return c.JSON(http.StatusOK, api.Result{
		Success: true,
	})
}
