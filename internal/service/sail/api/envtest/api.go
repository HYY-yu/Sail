// Package envtest 构造一个管理后台的集成测试环境
// 集成测试的接口封装，比如对这些配置文件的操作、测试环境的构建

package envtest

import (
	"context"
	"errors"
	"fmt"
	"github.com/HYY-yu/sail/internal/service/sail/api/envtest/model"
	"github.com/carlmjohnson/requests"
	"time"
)

type TestDataId struct {
	ProjectId      int `json:"project_id"`
	NamespaceId    int `json:"namespace_id"`
	PublicConfigId int `json:"public_config_id"`
	ConfigId       int `json:"config_id"`
	ConfigLinkId   int `json:"config_link_id"`
}

type ApiEnvTest struct {
	LoginToken string
	Data       *TestDataId
}

func (a *ApiEnvTest) UpdateTestConfig(isPublic bool, newValue string) error {
	configId := a.Data.ConfigId
	if isPublic {
		configId = a.Data.PublicConfigId
	}

	param := map[string]interface{}{
		"config_id": configId,
		"content":   newValue,
	}

	type ApiResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	resp := &ApiResponse{}
	err := requests.URL(model.TestAPIURL).
		Path(model.TestUpdateConfig).
		Bearer(a.LoginToken).
		BodyJSON(param).
		ToJSON(resp).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Message)
	}
	return nil
}

func (a *ApiEnvTest) MetaConfigYaml() (string, error) {
	type ApiResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}
	resp := &ApiResponse{}

	err := requests.URL(model.TestAPIURL).
		Path(model.TestMetaConfig).
		Bearer(a.LoginToken).
		ToJSON(resp).
		Param("temp", "K8S").
		ParamInt("project_id", a.Data.ProjectId).
		ParamInt("project_group_id", model.TestProjectGroupId).
		ParamInt("namespace_id", a.Data.NamespaceId).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(resp.Message)
	}

	return resp.Data, nil
}

func (a *ApiEnvTest) Start() error {
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
	a.LoginToken = token

	// 2. 调用接口 CreateData
	err = requests.URL(model.TestAPIURL).
		Path(model.TestCreateTestData).
		Bearer(token).
		CheckStatus(200).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}

	// 3. 缓存测试数据ID
	type ApiResponse struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    *TestDataId `json:"data"`
	}
	resp := &ApiResponse{}

	err = requests.URL(model.TestAPIURL).
		Path(model.TestGetTestData).
		Bearer(token).
		ToJSON(resp).
		CheckStatus(200).
		Fetch(getFiveSecondCtx())
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return errors.New(resp.Message)
	}

	a.Data = resp.Data
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
