# 微机原理资源池后端

## 项目介绍

本项目是微机原理资源池的后端服务，提供用户管理、资源管理、论坛交流和AI助手等功能的API接口。

## 技术栈

- Go语言 + Gin框架：提供RESTful API服务
- GORM：Go语言ORM库，用于数据库操作
- MySQL：关系型数据库，存储应用数据
- Minio：对象存储服务，用于存储上传的文件
- JWT：用于用户认证和授权
- Docker：容器化部署

## 功能模块

1. 用户管理：注册、登录、个人资料管理
2. 资源管理：上传、下载、搜索、分类管理
3. 积分系统：资源上传下载积分奖励和消费
4. 论坛交流：发帖、回复、分类讨论
5. AI助手：基于大语言模型的智能问答
6. 管理功能：资源审核、用户管理、数据统计

## 目录结构

```
├── config/         # 配置相关
├── controllers/    # 控制器
├── middleware/     # 中间件
├── migrations/     # 数据库迁移
├── models/         # 数据模型
├── routes/         # 路由定义
├── utils/          # 工具函数
├── .env            # 环境变量
├── Dockerfile      # Docker构建文件
├── go.mod          # Go模块定义
├── go.sum          # Go依赖校验
└── main.go         # 程序入口
```

## 新增功能

1. **数据库迁移**：使用GORM自动迁移功能，确保数据库结构与模型定义一致
2. **Minio文件存储工具**：封装Minio操作，提供文件上传、下载、URL生成等功能
3. **积分系统控制器**：管理用户积分的获取和消费，记录积分变动历史
4. **资源审核功能**：管理员审核用户上传的资源，确保内容质量
5. **错误处理中间件**：统一处理应用程序错误，提供友好的错误响应
6. **Docker配置**：提供容器化部署支持，简化部署流程

## 部署说明

### 本地开发

1. 克隆项目

```bash
git clone <项目地址>
cd backend
```

2. 安装依赖

```bash
go mod download
```

3. 配置环境变量（修改.env文件）

4. 运行项目

```bash
go run main.go
```

### Docker部署

使用docker-compose一键部署整个应用：

```bash
docker-compose up -d
```

这将启动后端API服务、MySQL数据库和Minio对象存储服务。

## API文档

### 用户相关

- `POST /api/register`：用户注册
- `POST /api/login`：用户登录
- `GET /api/user/profile`：获取用户资料
- `PUT /api/user/profile`：更新用户资料

### 资源相关

- `GET /api/resources`：获取资源列表
- `GET /api/resources/:id`：获取资源详情
- `POST /api/resources`：创建资源
- `PUT /api/resources/:id`：更新资源
- `DELETE /api/resources/:id`：删除资源
- `GET /api/resources/categories`：获取资源分类
- `GET /api/resources/search`：搜索资源
- `POST /api/upload`：上传文件
- `GET /api/download/:id`：下载文件

### 积分相关

- `GET /api/user/points`：获取用户积分
- `GET /api/user/points/history`：获取积分历史

### 管理员功能

- `GET /api/admin/resources/pending`：获取待审核资源
- `PUT /api/admin/resources/:id/review`：审核资源
- `POST /api/admin/points/add`：添加用户积分
- `GET /api/admin/stats`：获取统计信息

## 注意事项

1. 生产环境部署时，请修改JWT密钥和数据库密码等敏感信息
2. Minio服务需要确保存储空间足够
3. 定期备份数据库数据