# 开放式课程学习外延资源池系统

## 项目简介

本项目旨在打造一个开放式资源池系统，为用户提供一个集资源分享、学习交流、知识获取于一体的综合性平台。系统包含了用户管理、资源管理、论坛交流、在线聊天、积分激励以及后台管理等核心功能模块，致力于构建一个活跃、互助的学习社区。

## 系统架构

本系统采用前后端分离的架构：

- **后端 (Backend)**: 使用 Golang 语言，基于 Gin 框架开发，负责处理业务逻辑、数据存储与API接口服务。
- **前端 (Frontend)**: 使用 Vue.js 框架，结合 PrimeVue 组件库和 Tailwind CSS 进行界面开发，提供用户交互界面。
- **数据库 (Database)**: 主要使用 MySQL 存储核心业务数据。
- **缓存 (Cache)**: 使用 Redis 进行数据缓存，提升系统性能和响应速度。
- **对象存储 (Object Storage)**: 使用 MinIO 存储用户上传的资源文件，如文档、视频等。

## 主要功能模块

- **用户模块**: 用户注册、登录、个人资料管理、头像上传、JWT令牌刷新等。
- **资源模块**: 资源列表展示、分类浏览、资源搜索、资源详情查看、资源上传下载、评论、收藏、点赞/点踩等。
- **论坛模块**: 论坛板块、帖子发布与管理、回帖、帖子收藏、点赞等。
- **聊天模块**: (根据实际情况填写，例如：私聊、群聊、实时消息等)。
- **积分模块**: (根据实际情况填写，例如：用户行为积分、积分兑换等)。
- **管理模块**: (根据实际情况填写，例如：用户管理、资源审核、内容管理、系统配置等)。

## 技术栈

### 后端 (Golang)

- **Web框架**: Gin (https://gin-gonic.com/)
- **ORM**: GORM (https://gorm.io/)
- **数据库**: MySQL (https://www.mysql.com/)
- **缓存**: Redis (https://redis.io/)
- **对象存储**: MinIO (https://min.io/)
- **API认证**: JWT (JSON Web Tokens)
- **ID生成**: Snowflake

### 前端 (Vue.js)

- **核心框架**: Vue.js 3 (https://vuejs.org/)
- **UI组件库**: PrimeVue (https://primevue.org/)
- **CSS框架**: Tailwind CSS (https://tailwindcss.com/)
- **路由管理**: Vue Router (https://router.vuejs.org/)
- **状态管理**: Pinia (https://pinia.vuejs.org/)
- **构建工具**: Vite (https://vitejs.dev/)
- **HTTP客户端**: Axios (https://axios-http.com/)

## API 文档

详细的API接口文档请参考：[API 文档](./backend/docs/api_documentation.md)

## 安装与启动指南

### 后端 (Golang)

1.  **环境准备**:
    *   安装 Go (推荐版本 1.18+)
    *   安装 MySQL
    *   安装 Redis
    *   安装 MinIO (或使用云服务商提供的对象存储)

2.  **克隆项目**:
    ```bash
    git clone <your-repository-url>
    cd V1-BK/backend
    ```

3.  **配置环境变量**:
    *   复制 `.env.example` 文件为 `.env`。
    *   根据您的实际环境修改 `.env` 文件中的数据库连接信息、Redis连接信息、MinIO配置、JWT密钥等。

4.  **安装依赖**:
    ```bash
    go mod tidy
    ```

5.  **数据库迁移** (如果项目包含数据库迁移脚本):
    ```bash
    # 根据项目实际情况执行迁移命令，例如：
    # go run main.go migrate (假设迁移命令集成在main.go中)
    # 或者手动执行SQL迁移脚本
    ```
    请确保 `migrations/migrate.go` 中的逻辑能够正确初始化数据库表结构。

6.  **启动服务**:
    ```bash
    go run main.go
    ```
    后端服务默认运行在 `http://localhost:8080` (或其他在配置中指定的端口)。

### 前端 (Vue.js)

1.  **环境准备**:
    *   安装 Node.js (推荐 LTS 版本)
    *   安装 pnpm (推荐) 或 npm/yarn

2.  **进入前端目录**:
    ```bash
    cd ../frontend 
    # 或者从项目根目录 cd V1-BK/frontend
    ```

3.  **配置环境变量** (如果需要):
    *   检查是否存在 `.env` 或 `.env.development` 文件，并根据需要配置API基础路径等。
    *   通常API基础路径会指向后端服务的地址，例如 `VITE_API_BASE_URL=http://localhost:8080/api`。

4.  **安装依赖**:
    ```bash
    # 使用 pnpm
    pnpm install
    
    # 或者使用 npm
    # npm install
    
    # 或者使用 yarn
    # yarn install
    ```

5.  **启动开发服务器**:
    ```bash
    # 使用 pnpm
    pnpm dev
    
    # 或者使用 npm
    # npm run dev
    
    # 或者使用 yarn
    # yarn dev
    ```
    前端开发服务器通常运行在 `http://localhost:5173` (Vite 默认端口) 或其他指定端口。

6.  **构建生产版本**:
    ```bash
    # 使用 pnpm
    pnpm build
    
    # 或者使用 npm
    # npm run build
    
    # 或者使用 yarn
    # yarn build
    ```
    构建后的静态文件会输出到 `dist` 目录。

