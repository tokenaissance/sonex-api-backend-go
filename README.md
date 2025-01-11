# sonex-api-backend-go

## development

基于 [go-zero](https://go-zero.dev/) 框架。


### 新增 API

步骤：

在 `api/` 中新增 API 文件，并在 `sonex_api_backend_go.api` 中 `import`.

重新生成 API 文件对应的 Go 代码：

```sh
goctl api go --api sonex_api_backend_go.api --style go_zero --dir .
```

重新生成 Swagger 文件：

```sh
goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api sonex_api_backend_go.api -dir .
mv swagger.json ./static/swagger
```

更新 model:

```sh
goctl model pg datasource --url="postgres://postgres:postgres@10.0.0.150:5432/sonex_api_backend?sslmode=disable" --dir ./internal/repository --style go_zero --table <表名> 
```
