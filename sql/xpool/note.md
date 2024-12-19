
## 生成api
goctl api go --api sonex_api_backend_go.api --style go_zero --dir .
goctl api validate --api sonex_api_backend_go.api

## 生成model
goctl model pg datasource --url="postgres://postgres:123456Abc@localhost:5432/sonex_api_backend?sslmode=disable" --dir ./internal/repository --style go_zero --table "xpool_pools,xpool_user_action"

## 生成swagger
goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api sonex_api_backend_go.api -dir . && mv swagger.json ./static/swagger


