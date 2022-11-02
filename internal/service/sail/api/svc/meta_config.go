package svc

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	"github.com/HYY-yu/seckill.pkg/core"
	"github.com/HYY-yu/seckill.pkg/pkg/response"

	"github.com/HYY-yu/sail/internal/service/sail/config"
	"github.com/HYY-yu/sail/internal/service/sail/model"
)

type MetaConfig struct {
	ETCDEndpoints  string // 逗号分隔的ETCD地址，0.0.0.0:2379,0.0.0.0:12379,0.0.0.0:22379
	ETCDUsername   string
	ETCDPassword   string
	ProjectKey     string
	Namespace      string
	NamespaceKey   string
	Configs        string // 逗号分隔的 config_name.config_type，如：mysql.toml,cfg.json,redis.yaml，空代表不下载任何配置
	ConfigFilePath string // 本地配置文件存放路径，空代表不存储本都配置文件
	LogLevel       string // 日志级别(DEBUG\INFO\WARN\ERROR)，默认 WARN
	MergeConfig    bool   // 是否合并配置，合并配置则会将同类型的配置合并到一个文件中，需要先设置ConfigFilePath
}

const flagTemplate = `
--sail-etcd-endpoints={{.ETCDEndpoints}} \
{{if .ETCDUsername}}--sail-etcd-username={{.ETCDUsername}}{{end}} {{if .ETCDPassword}}--sail-etcd-password={{.ETCDPassword}}{{end}} \
--sail-project-key={{.ProjectKey}} --sail-namespace={{.Namespace}} \
{{if .NamespaceKey}}--sail-namespace-key={{.NamespaceKey}}{{end}} {{if .Configs}}--sail-configs={{.Configs}}{{end}} \ 
{{if .ConfigFilePath}}--sail-config-file-path={{.ConfigFilePath}}{{end}} {{if .LogLevel}}--sail-log-level={{.LogLevel}}{{end}} {{if .MergeConfig}}--sail-merge-config={{.MergeConfig}}{{end}}
`

const envTemplate = `
export SAIL_ETCD_ENDPOINTS={{.ETCDEndpoints}}
{{if .ETCDUsername}}
export SAIL_ETCD_USERNAME={{.ETCDUsername}}
{{end}}
{{if .ETCDPassword}}
export SAIL_ETCD_PASSWORD={{.ETCDPassword}} 
{{end}}
export SAIL_PROJECT_KEY={{.ProjectKey}}
export SAIL_NAMESPACE={{.Namespace}}
{{if .NamespaceKey}}
export SAIL_NAMESPACE_KEY={{.NamespaceKey}}
{{end}}
{{if .Configs}}
export SAIL_CONFIGS={{.Configs}} 
{{end}}
{{if .ConfigFilePath}}
export SAIL_CONFIG_FILE_PATH={{.ConfigFilePath}}
{{end}}
{{if .LogLevel}}
export SAIL_LOG_LEVEL={{.LogLevel}}
{{end}}
{{if .MergeConfig}}
export SAIL_MERGE_CONFIG={{.MergeConfig}}
{{end}}
`

const tomlTemplate = `
[sail]
etcd_endpoints = "{{.ETCDEndpoints}}"
{{if .ETCDUsername}}
etcd_username="{{.ETCDUsername}}"
{{end}}
{{if .ETCDPassword}}
etcd_password="{{.ETCDPassword}}"
{{end}}
project_key = "{{.ProjectKey}}"
namespace = "{{.Namespace}}"
{{if .NamespaceKey}}
namespace_key="{{.NamespaceKey}}"
{{end}}
{{if .Configs}}
configs="{{.Configs}}"
{{end}}
{{if .ConfigFilePath}}
config_file_path="{{.ConfigFilePath}}"
{{end}}
{{if .LogLevel}}
log_level="{{.LogLevel}}"
{{end}}
{{if .MergeConfig}}
merge_config={{.MergeConfig}}
{{end}}
`

func getTemplateStringBy(temp string, mc *MetaConfig) (string, error) {
	t := template.Must(template.New("template").Parse(temp))
	sw := bytes.Buffer{}
	err := t.Execute(&sw, mc)
	return sw.String(), err
}

func (s *ConfigSvc) GetTemplate(sctx core.SvcContext, temp string, projectID int, projectGroupID int, namespaceID int) (string, error) {
	ctx := sctx.Context()
	mgr := s.ConfigRepo.Mgr(ctx, s.DB.GetDb())
	mgr.WithPrepareStmt()

	project, namespace, err := s.getConfigProjectAndNamespace(ctx, projectID, namespaceID)
	if err != nil {
		return "", response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	if len(namespace.SecretKey) > 0 {
		// 如果有秘钥，则要求 Owner 才能访问
		_, role := s.CheckStaffGroup(ctx, projectGroupID)
		if role > model.RoleOwner {
			return "", response.NewErrorWithStatusOk(
				response.AuthorizationError,
				"此命名空间您无权访问",
			)
		}
	}

	metaConfig := &MetaConfig{
		ETCDEndpoints:  strings.Join(config.Get().ETCD.Endpoints, ","),
		ETCDUsername:   config.Get().ETCD.Username,
		ETCDPassword:   config.Get().ETCD.Password,
		ProjectKey:     project.Key,
		Namespace:      namespace.Name,
		NamespaceKey:   namespace.SecretKey,
		ConfigFilePath: config.Get().SDK.ConfigFilePath,
		LogLevel:       config.Get().SDK.LogLevel,
		MergeConfig:    config.Get().SDK.MergeConfig,
	}
	configList, err := mgr.WithOptions(
		mgr.WithProjectID(projectID),
		mgr.WithProjectGroupID(projectGroupID),
		mgr.WithNamespaceID(namespaceID),
	).WithSelects(
		model.ConfigColumns.ID,
		model.ConfigColumns.Name,
		model.ConfigColumns.NamespaceID,
		model.ConfigColumns.ConfigType,
	).
		Gets()
	if err != nil {
		return "", response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	configs := make([]string, len(configList))
	for i, m := range configList {
		configs[i] = m.Name + "." + m.ConfigType
	}

	metaConfig.Configs = strings.Join(configs, ",")

	tempStr := ""
	switch temp {
	case "FLAG":
		tempStr = flagTemplate
	case "ENV":
		tempStr = envTemplate
	case "TOML":
		tempStr = tomlTemplate
	default:
		return "", nil
	}
	result, err := getTemplateStringBy(tempStr, metaConfig)
	if err != nil {
		return "", response.NewErrorAutoMsg(
			http.StatusInternalServerError,
			response.ServerError,
		).WithErr(err)
	}
	return result, nil
}
