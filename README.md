# GeoNames Service

一个基于 Go 语言的地理位置数据服务，提供地理位置数据的下载和查询功能。

## 安装指南

1. 克隆仓库:
```bash
git clone https://github.com/your-repo/geonames-service.git
```
2. 安装依赖:
```bash
go mod download
```
3. 配置数据库:
   - 创建PostgreSQL数据库
   - 修改config.yaml中的数据库连接信息
4. 运行迁移:
```bash
go run cmd/cli/main.go migrate
```
5. 下载地理位置数据并导入数据库:
```bash
go run cmd/cli/main.go download
```
6. 启动服务:
```bash
go run cmd/server/main.go
```

## 功能特性

- 支持下载地理位置数据
- 提供地理位置数据查询接口
- 支持按国家代码筛选地理位置信息

## 技术栈

- Go 1.23.3
- Gorilla Mux (HTTP 路由)
- PostgreSQL (数据存储)


## API 接口


### 获取地理位置列表

```
GET /locations
```

获取所有地理位置数据（默认返回前100条记录）。

### 按国家代码查询

```
GET /locations/{countryCode}
```

根据国家代码查询地理位置数据。

请求示例:
```bash
curl http://localhost:8080/locations/CN
```

响应:
```json
[
  {
    "id": 1,
    "name": "Beijing",
    "countryCode": "CN",
    "latitude": 39.9042,
    "longitude": 116.4074
  }
]
```

## 许可证

MIT License