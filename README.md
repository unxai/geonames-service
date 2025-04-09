# GeoNames Service

一个基于 Go 语言的地理位置数据服务，提供地理位置数据的下载和查询功能。

## 功能特性

- 支持下载地理位置数据
- 提供地理位置数据查询接口
- 支持按国家代码筛选地理位置信息

## 技术栈

- Go 1.23.3
- Gorilla Mux (HTTP 路由)
- PostgreSQL (数据存储)

## API 接口

### 下载数据

```
POST /download
```

下载并更新地理位置数据。

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

## 许可证

MIT License