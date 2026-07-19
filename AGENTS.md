# Sail Agent Guide

## 项目概览

Sail 是 Go 编写的配置中心，包含两个独立 Go 模块：

- 根模块（Go 1.20）：主服务，入口为 `internal/service/sail/cmd/main.go`。
- `internal/operator`（Go 1.21）：Kubernetes Operator，拥有独立的 `go.mod`、Makefile、CRD 和部署配置。

主要目录：

- `internal/service/sail/api/`：HTTP API、Handler、Service、Repository、依赖注入。
- `internal/service/sail/storage/`：etcd 存储实现。
- `internal/service/sail/config/`：主服务配置。
- `internal/service/sail/docs/`：Swagger 生成产物。
- `ui/`：内嵌的静态资源和 HTML 模板。
- `sql_backup/db.sql`：本地/初始环境数据库建表脚本。
- `internal/operator/`：ConfigMapRequest Kubernetes Controller。
- `deploy/`、`docker-compose.yml`：容器化和本地依赖环境。

## 工作原则

- 先定位完整调用链，再修改：路由 → Handler → Service → Repository/Storage。
- 改动保持最小范围；涉及权限、配置发布、加密、etcd 同步或数据库结构时，先说明影响范围。
- 不修改生成文件；修改其源定义后，使用对应生成命令更新产物。
- 不提交密钥、JWT、密码、Namespace Key、解密后的配置内容或本地环境配置。
- 关键业务边界应记录结构化 `info` 日志：操作对象标识、请求结果、耗时和错误原因。严禁记录密码、Token、密钥或完整配置正文。
- 未经明确授权，不执行 Docker push、Kubernetes `apply/delete/deploy`、生产配置发布或破坏性数据库操作。

## 根服务开发

在仓库根目录执行：

```bash
go mod download
go build -o ./tmp/sail ./internal/service/sail/cmd/main.go
go vet ./...
```

本地运行：

```bash
go run ./internal/service/sail/cmd/main.go
```

服务默认从 `internal/service/sail/config/cfg.toml` 读取本地配置；容器环境可使用：

```bash
go run ./internal/service/sail/cmd/main.go -config_path <cfg.toml 所在目录>
```

本地依赖环境：

```bash
docker compose up -d mysql etcd-node
```

根模块测试依赖本地 etcd，默认使用 `127.0.0.1:2379`；运行完整测试前确认 etcd 已启动：

```bash
go test -v ./...
```

测试会写入并清理测试用 etcd Key。不要将测试指向共享或生产 etcd。

## API、数据与配置约束

- 新增或变更 API 时，同时检查路由分组、中间件和权限边界；不要意外绕过 JWT 或员工组鉴权。
- 修改公开 API、DTO 或注释后，更新 Swagger 生成产物。
- `sql_backup/db.sql` 是初始化脚本。对已部署环境的结构变更应新增可追踪、可回滚或幂等的迁移脚本，不要把修改历史仅留在初始化 DDL 中。
- 配置中心存储敏感业务配置：排查时使用配置 ID、项目 Key、命名空间和版本号，不输出原始配置或解密结果。
- UI 修改应保持与 `/ui` 路由和现有模板结构一致；本项目没有独立的前端构建流程。

## 生成文件

以下文件或目录为生成产物，禁止直接手改：

- `internal/service/sail/docs/docs.go`
- `internal/service/sail/api/wire_gen.go`
- `internal/service/sail/api/repo/gen*.go`
- `internal/service/sail/api/repo/sail.gen.*.go`
- `internal/service/sail/api/repo/*_mock.go`
- `internal/service/sail/api/svc_interface/mock.go`
- `internal/operator/api/v1beta1/zz_generated.deepcopy.go`
- `internal/operator/config/crd/` 下的 CRD 产物

按需生成，避免在无关改动中批量重写文件：

```bash
# Swagger
cd internal/service/sail
go generate pkg.go

# Wire 依赖注入
cd internal/service/sail/api
go generate wire.go

# GORMT Repository：需要本地 Sail MySQL schema，运行前先确认数据库和影响范围
cd internal/service/sail/api/repo
go generate gen.go
```

## Kubernetes Operator

所有 Operator 命令必须在 `internal/operator` 中执行：

```bash
cd internal/operator
make build
make test
make lint
```

- 修改 `api/v1beta1` 的 CRD 类型后，运行 `make manifests generate`。
- `make test` 会下载/使用 envtest 资产，并可能更新生成文件或格式化代码。
- `make install`、`make deploy`、`make uninstall`、`make undeploy` 会操作当前 kubeconfig 指向的集群；执行前必须获得明确授权，并先确认当前 Kubernetes context。
- 修改 Controller 时，关注 Reconcile 的幂等性、错误重试和 Status 更新，避免在循环中产生无意义写入或日志噪声。

## 提交前检查

根据改动范围执行最小且充分的验证：

```bash
# 修改主服务 Go 代码
gofmt -w <changed-go-files>
go vet ./...
go test ./...

# 修改 Operator
cd internal/operator
make test
```

提交前检查：

```bash
git status --short
git diff --check
```

不要混入无关格式化、自动生成文件重写、IDE 文件、构建产物或本地配置。
