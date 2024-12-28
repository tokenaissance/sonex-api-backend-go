# sonex-api-backend-go

### 新增/修改API接口
* 根据 sonex_api_backend_go.api import 路径在 /api 目录中新建/修改 <group>.api 文件
* 在对应 api 文件中定义 Request/Reponse/Service 等结构
* 在根目录运行 goctl api go --api sonex_api_backend_go.api --style go_zero --dir . 生成接口相关文件
* 在根目录运行 goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api sonex_api_backend_go.api -dir . && mv swagger.json ./static/swagger 更新 swagger 信息