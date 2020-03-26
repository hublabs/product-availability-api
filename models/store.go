package models

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hublabs/product-availability-api/config"
)

type Store struct {
	Id             int64               `json:"id" xorm:"int64 notnull autoincr pk 'id'"`
	TenantId       int64               `json:"tenantId"  xorm:"unique(code) index notnull"`
	Code           string              `json:"code" xorm:" index notnull"`
	Name           string              `json:"name" xorm:"notnull"`
	Manager        string              `json:"manager" xorm:"notnull"`
	TelNo          string              `json:"telNo"`
	Area           string              `json:"area" xorm:"notnull"`
	Address        string              `json:"address" xorm:"notnull"`
	StatusCode     string              `json:"statusCode" xorm:"notnull"`
	Cashier        bool                `json:"cashier" xorm:"notnull"`
	ContractNo     string              `json:"contractNo"`
	OpenDate       string              `json:"openDate" xorm:"notnull"`
	CloseDate      string              `json:"closeDate" xorm:"notnull"`
	Remark         Remark              `json:"remark" xorm:"json"`
	Enable         bool                `json:"enable" xorm:"notnull"`
	CreatedAt      time.Time           `json:"-" xorm:"created"`
	UpdatedAt      time.Time           `json:"-" xorm:"updated"`
	Version        int                 `json:"-" xorm:"version"`
	PaymentMethods []PaymentMethodView `json:"paymentMethods" xorm:"-"`
	SaleChannels   []SaleChannelView   `json:"saleChannels" xorm:"-"`
	Brands         []BrandView         `json:"brands" xorm:"-"`
}

type PaymentMethodView struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	EId    int64  `json:"eId"`
	Enable bool   `json:"enable"`
}

type SaleChannelView struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type BrandView struct {
	Id     int64  `json:"id"`
	Code   string `json:"code"`
	Enable bool   `json:"enable"`
}

type Remark struct {
	ElandShopInfos []ElandShopInfo `json:"elandShopInfos"`
}

type ElandShopInfo struct {
	BrandCode string `json:"brandCode"`
	ShopCode  string `json:"shopCode"`
}

func (Store) getById(ctx context.Context, storeId int64) (Store, error) {
	var resp struct {
		Result struct {
			TotalCount int     `json:"totalCount"`
			Items      []Store `json:"items"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}

	if err := RetryWebApi(ctx, &resp, http.MethodGet, config.Config().Services.PlaceManagementApi+"/v1/store/getallinfo?storeids="+strconv.FormatInt(storeId, 10)+"&MaxResultCount=1", nil); err != nil {
		return Store{}, err
	}

	if !resp.Success {
		return Store{}, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result.Items[0], nil
}
