# Gin Naive Admin（后端）

Gin 实现的 REST API 服务，为 **Gin Naive Admin** 前端提供登录鉴权、用户/角色/菜单、部门、字典、职务、文件上传、仪表盘、操作与登录日志、系统参数等能力。

## 技术栈

- **Go**（见 `go.mod` 中版本要求）
- **Gin** — HTTP 框架
- **GORM** + **MySQL** — 持久化
- **Redis** — 缓存等
- **JWT** — 访问令牌与刷新
- **Viper** — YAML 配置与热更新
- **Zap** — 结构化日志（支持文件轮转）
- **验证码**（base64Captcha）等
- **Swagger**（[swaggo/swag](https://github.com/swaggo/swag) + gin-swagger）— OpenAPI 2.0 与 `/swagger` UI

## 工程结构（简要）

| 路径 | 说明 |
|------|------|
| `main.go` | 入口：配置、日志、校验器、数据库、Redis、路由 |
| `initialize/` | 初始化：viper、gorm、redis、路由、中间件、静态资源 |
| `api/` | 按模块划分的路由与业务（`user`、`role`、`menu`、`login` 等） |
| `middleware/` | JWT、CORS、操作日志、安全头等 |
| `model/` | 数据模型与响应结构 |
| `permission/` | 权限注册、菜单/按钮种子、用户有效权限解析等 |
| `global/` | 全局配置与单例（DB、Redis、Logger 等） |
| `config.development.yaml` / `config.release.yaml` | 环境与运行配置示例 |
| `docs/` | `swag` 生成的 OpenAPI 文档（`docs.go`、`swagger.json`、`swagger.yaml`） |

对外 API 统一前缀由配置项 **`router.router-prefix`** 控制，默认为 **`/api`**（例如健康检查：`GET /api/health`）。

## 环境依赖

- **MySQL**（`utf8mb4`）  
- **Redis**

请在配置文件中填写正确的连接信息，并提前创建数据库（库名与 `database.db_name` 一致）。

## 配置说明

使用 **YAML** 配置文件，通过启动参数 **`-r`** 指定文件路径；也可通过环境变量 **`GNA_CONFIG`** 指定配置文件路径。

常用配置块：

- **`system`**：环境、`port`、应用名  
- **`database`**：MySQL 连接池与库名  
- **`redis`**：地址、库索引、密码  
- **`jwt`**：密钥、过期与刷新时间（**生产环境务必修改密钥**）  
- **`router`**：`router-prefix`、静态上传目录 `path`（如 `uploads/file`）  
- **`security`**：`super-role-codes`（超级角色编码，如 `admin`）、数据范围等  
- **`cors`**：跨域白名单或放行策略（需包含前端开发/生产域名）  
- **`zap`**：日志级别、目录、轮转  

首次启动会在连接数据库后执行表注册及部分种子数据（菜单按钮权限、系统默认配置等），具体逻辑见 `main.go` 与 `permission`、`sysconfig` 包。

## OpenAPI / Swagger 文档

服务启动后，在浏览器打开：

- **Swagger UI**：`http://127.0.0.1:<port>/swagger/index.html`（端口见 `system.port`，默认 `8080`）

说明：

- 文档中的 **`BasePath`** 为 `/api`，与配置项 `router.router-prefix` 一致；若修改前缀，请在 `main.go` 顶部的 swag 注释中同步修改 `@BasePath` 并重新生成文档。
- 需登录的接口在 Swagger 中已标记 **`AccessToken`**，调试前点击右上角 **Authorize**，填入登录后获得的 JWT（请求头字段名 **`AccessToken`**，与前端一致）。
- 修改或新增接口注释后，在项目根目录执行以下命令重新生成 `docs/`：

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.4 init -g main.go -o docs --parseDependency --parseInternal
```

## 运行方式

在项目根目录（本 `README.md` 所在目录）执行：

```bash
# 显式指定配置文件（推荐）
go run main.go -r config.development.yaml
```

未传 `-r` 时，会按 Gin 运行模式选择默认文件（开发模式下多为 `config.development.yaml`），也可设置环境变量 `GNA_CONFIG` 指向配置文件。

服务监听地址为配置中的 `system.port`，控制台会打印本机访问 URL（如 `http://127.0.0.1:8080`）。

### 构建二进制

```bash
go build -o gin-naive-admin .
./gin-naive-admin -r config.release.yaml
```

## 与前端联调

1. 启动本服务并保证 MySQL、Redis 可用。  
2. 在 `cors.whitelist` 中加入前端开发地址（例如 `http://127.0.0.1:5173`）。  
3. 前端 `.env.development` 中 `VITE_API_BASE_URL` 应指向本服务根地址（如 `http://127.0.0.1:8080`），请求路径为 `/api/...`。

详见同级目录 [gin-naive-admin-web/README.md](../gin-naive-admin-web/README.md)。

## 补充脚本

`scripts/` 下为历史或增量 SQL 片段，可按迁移需要参考使用；以应用内 GORM 迁移与种子逻辑为准。

## 安全提示

- 勿将含真实密码、JWT 密钥的配置文件提交到公开仓库。  
- 生产部署请使用 `config.release.yaml` 或等价配置，并关闭调试、收紧 CORS。
