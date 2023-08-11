// Package envtest 构造一个管理后台的集成测试环境
// 集成测试的接口封装，比如对这些配置文件的操作、测试环境的构建

package envtest

import (
	"context"
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/api/envtest/model"
	"github.com/carlmjohnson/requests"
	"time"
)

type ApiEnvTest struct {
}

func (*ApiEnvTest) Start() error {
	err := requests.URL(model.TestAPIURL).
		Path(model.TestCheckHeath).
		CheckStatus(200).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}
	// 1. 登录
	token, err := getLoginAccessToken()
	if err != nil {
		return err
	}

	// 2. 调用接口 CreateData
	err = requests.URL(model.TestAPIURL).
		Path(model.TestCreateTestData).
		Bearer(token).
		CheckStatus(200).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}
	return nil
}

func (*ApiEnvTest) Stop() error {
	token, err := getLoginAccessToken()
	if err != nil {
		return err
	}

	// 调用 CleanData 接口清除测试数据
	err = requests.URL(model.TestAPIURL).
		Path(model.TestCleanTestData).
		Bearer(token).
		CheckStatus(200).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}
	return nil
}

func getFiveSecondCtx() context.Context {
	timeCtx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return timeCtx
}

func getLoginAccessToken() (string, error) {
	loginReq := map[string]interface{}{
		"user_name": model.TestAccountName,
		"password":  model.TestAccountPass,
	}

	loginResp := struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
	}{}
	err := requests.URL(model.TestAPIURL).
		Path(model.TestLogin).
		BodyJSON(loginReq).
		CheckStatus(200).
		ToJSON(&loginResp).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return "", err
	}

	if loginResp.Code != 0 {
		return "", fmt.Errorf("login err: %s", loginResp.Message)
	}
	return loginResp.Data.AccessToken, nil
}
