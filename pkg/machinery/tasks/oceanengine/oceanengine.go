package oceanengine

import (
	"context"
	"fmt"

	"github.com/oceanengine/ad_open_sdk_go"
	oc_config "github.com/oceanengine/ad_open_sdk_go/config"
	"github.com/oceanengine/ad_open_sdk_go/models"
)

type OceanEngineOpenSdk struct {
	client *ad_open_sdk_go.Client
}

func NewOceanEngineOpenSdk() *OceanEngineOpenSdk {
	client := ad_open_sdk_go.Init(oc_config.NewConfiguration())
	return &OceanEngineOpenSdk{client: client}
}

// GetAccessToken 获取access_token
func (o *OceanEngineOpenSdk) GetAccessToken(ctx context.Context, appId *int64, secret, authCode string) (at *models.Oauth2AccessTokenResponseData, err error) {
	var req models.Oauth2AccessTokenRequest

	req.AppId = appId
	req.Secret = secret
	req.AuthCode = authCode
	resp, _, err := o.client.Oauth2AccessTokenApi().
		Post(ctx).
		Oauth2AccessTokenRequest(req).
		Execute()

	if err != nil {
		return nil, err
	}

	if *resp.Code != 0 {
		return nil, fmt.Errorf("code:%d,message:%s", *resp.Code, *resp.Message)
	}

	at = resp.Data

	return
}

// RefreshAccessToken 刷新access_token
func (o *OceanEngineOpenSdk) RefreshAccessToken(ctx context.Context, appId *int64, secret, refreshToken string) (at *models.Oauth2RefreshTokenResponseData, err error) {
	var req models.Oauth2RefreshTokenRequest

	req.AppId = appId
	req.Secret = secret
	req.RefreshToken = refreshToken

	resp, _, err := o.client.Oauth2RefreshTokenApi().
		Post(ctx).
		Oauth2RefreshTokenRequest(req).
		Execute()

	if err != nil {
		return nil, err
	}

	if *resp.Code != 0 {
		return nil, fmt.Errorf("code:%d,message:%s", *resp.Code, *resp.Message)
	}

	return resp.Data, nil
}

type CustomerCenterAdvertiserListGetRequest struct {
	AccountSource *models.CustomerCenterAdvertiserListV2AccountSource
	AccessToken   string
	CcAccountId   int64
	Filtering     models.CustomerCenterAdvertiserListV2Filtering
	Page          int64
	PageSize      int64
}

// GetAdvertiserIdList 获取广告主id列表
func (o *OceanEngineOpenSdk) GetAdvertiserIdList(ctx context.Context, req CustomerCenterAdvertiserListGetRequest) (result *models.CustomerCenterAdvertiserListV2ResponseData, err error) {
	res, _, err := o.client.CustomerCenterAdvertiserListV2Api().
		Get(ctx).
		AccessToken(req.AccessToken).
		CcAccountId(req.CcAccountId).
		Filtering(req.Filtering).
		Page(req.Page).
		PageSize(req.PageSize).
		Execute()
	if err != nil {
		return nil, err
	}

	if *res.Code != 0 {
		return nil, fmt.Errorf("code:%d,message:%s", *res.Code, *res.Message)
	}

	return res.Data, err
}

type ReportAdvertiserGetRequest struct {
	AccessToken     string
	AdvertiserId    int64
	StartDate       string
	EndDate         string
	Fields          []string
	TimeGranularity models.ReportAdvertiserGetV2TimeGranularity
	Filtering       models.ReportAdvertiserGetV2Filtering
}

// ReportAdvertiserGet 获取报表数据
func (o *OceanEngineOpenSdk) ReportAdvertiserGet(ctx context.Context, req *ReportAdvertiserGetRequest) (result []models.ReportAdvertiserGetV2Response, err error) {
	resp, _, err := o.client.ReportAdvertiserGetV2Api().
		Get(ctx).
		AccessToken(req.AccessToken).
		AdvertiserId(req.AdvertiserId). // 广告主ID
		StartDate(&req.StartDate).      // 开始日期
		EndDate(&req.EndDate).          // 结束日期
		Fields(req.Fields).             // 字段
		Filtering(req.Filtering).       // 过滤
		// GroupBy(request.GroupBy).                 // 分组
		// OrderField(request.OrderField).           // 排序字段
		// OrderType(request.OrderType).             // 排序类型
		// Page(request.Page).                       // 页码
		PageSize(50).                         // 每页条数
		TimeGranularity(req.TimeGranularity). // 时间粒度
		Execute()

	if err != nil {
		return nil, err
	}

	if *resp.Code != 0 {
		return nil, fmt.Errorf("code:%d,message:%s", *resp.Code, *resp.Message)
	}

	return
}
