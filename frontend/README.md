# 微机原理资源池前端

这是开放式微机原理资源池项目的前端，基于Vue 3、Vite和TailwindCSS构建，提供了微机原理相关资源管理、分享和下载等功能的现代化Web界面。

## 功能特点

- **用户认证**：注册、登录和个人资料管理
- **资源中心**：浏览、搜索、上传和下载资源
- **论坛交流**：话题讨论和回复功能
- **AI聊天**：智能AI助手对话
- **积分系统**：资源下载和上传的积分管理

## 技术栈

- **Vue 3**：渐进式JavaScript框架
- **Vue Router**：官方路由管理器
- **Pinia**：Vue的状态管理库
- **Vite**：现代前端构建工具
- **TailwindCSS**：实用优先的CSS框架
- **PrimeVue**：Vue UI组件库
- **Vuelidate**：表单验证库
- **Axios**：基于Promise的HTTP客户端

## 项目结构

```
/frontend
├── public/             # 静态资源
├── src/
│   ├── api/            # API服务
│   ├── assets/         # 资源文件
│   ├── components/     # 公共组件
│   ├── router/         # 路由配置
│   ├── stores/         # Pinia状态管理
│   ├── views/          # 页面组件
│   ├── App.vue         # 根组件
│   └── main.js         # 入口文件
├── index.html          # HTML模板
├── package.json        # 项目依赖
├── vite.config.js      # Vite配置
└── tailwind.config.js  # Tailwind配置
```

## 开发指南

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 后端API

前端通过RESTful API与后端通信，API基础路径为`/api`。详细的API文档请参考后端项目的说明文档。

## 浏览器兼容性

支持所有现代浏览器，包括：

- Chrome
- Firefox
- Safari
- Edge

## 许可证

[MIT](LICENSE)