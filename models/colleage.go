package models

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hublabs/product-availability-api/config"
)

type Colleague struct {
	Id          int64  `json:"id"`
	ColleagueNo string `json:"colleagueNo"`
	UserName    string `json:"userName"`
	Mobile      string `json:"mobile"`
	Unionid     string `json:"unionid"`
	Name        string `json:"name"`
}

func (Colleague) GetById(ctx context.Context, colleagueId int64) (*Colleague, error) {
	var resp struct {
		Result struct {
			Colleague Colleague `json:"colleague"`
		} `json:"result"`
		Success bool `json:"success"`
		Error   struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}

	url := fmt.Sprintf("%s/api/v1/login/colleague/login-token/%d?appCode=CloudPortal", config.Config().Services.ColleagueAuthApi, colleagueId)

	if err := RetryWebApi(ctx, &resp, http.MethodGet, url, nil); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("%d-%s", resp.Error.Code, resp.Error.Message)
	}

	return &resp.Result.Colleague, nil
}
