# API 文档

## 用户认证 (Authentication)
### 用户注册
- **方法**: POST
- **路径**: /register
- **请求体**:
  ```json
  {
    "username": "用户名",
    "email": "邮箱",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "user": {
      "id": 用户ID,
      "username": "用户名",
      "email": "邮箱",
      "points": 积分,
      "role": "角色"
    },
    "token": "JWT令牌"
  }
  ```

### 用户登录
- **方法**: POST
- **路径**: /login
- **请求体**:
  ```json
  {
    "username": "用户名或邮箱",
    "password": "密码"
  }
  ```
- **响应**:
  ```json
  {
    "user": {
      "id": 用户ID,
      "username": "用户名",
      "email": "邮箱",
      "avatar": "头像URL",
      "points": 积分,
      "role": "角色"
    },
    "token": "JWT令牌"
  }
  ```

### 获取用户资料
- **方法**: GET
- **路径**: /user/profile
- **请求头**:
  - Authorization: Bearer <token>
- **响应**:
  ```json
  {
    "id": 用户ID,
    "username": "用户名",
    "email": "邮箱",
    "avatar": "头像URL",
    "points": 积分,
    "role": "角色"
  }
  ```

### 更新用户资料
- **方法**: PUT
- **路径**: /user/profile
- **请求头**:
  - Authorization: Bearer <token>
- **请求体**:
  ```json
  {
    "email": "新邮箱",
    "avatar": "新头像URL"
  }
  ```
- **响应**:
  ```json
  {
    "id": 用户ID,
    "username": "用户名",
    "email": "邮箱",
    "avatar": "头像URL",
    "points": 积分,
    "role": "角色"
  }
  ```

### 获取用户积分
- **方法**: GET
- **路径**: /user/points
- **请求头**:
  - Authorization: Bearer <token>
- **响应**:
  ```json
  {
    "points": 积分
  }
  ```

## 资源控制器 (ResourceController)

### 获取资源列表
- **方法**: GET
- **路径**: /resources
- **参数**:
  - page (可选): 页码，默认1
  - pageSize (可选): 每页数量，默认10
  - category (可选): 分类ID
  - sort (可选): 排序方式，可选newest/popular，默认newest
- **响应**:
  ```json
  {
    "resources": [资源数组],
    "total": 总数,
    "page": 当前页,
    "pageSize": 每页数量
  }
  ```

### 获取资源详情
- **方法**: GET
- **路径**: /resources/:id
- **参数**:
  - id (路径参数): 资源ID
- **响应**:
  ```json
  {
    "id": 资源ID,
    "title": "标题",
    "description": "描述",
    "file_path": "文件路径",
    "status": "状态",
    "created_at": "创建时间",
    "updated_at": "更新时间",
    "User": {
      "id": 用户ID,
      "username": "用户名"
    },
    "Category": {
      "id": 分类ID,
      "name": "分类名称"
    }
  }
  ```

### 获取资源分类
- **方法**: GET
- **路径**: /resources/categories
- **响应**:
  ```json
  [
    {
      "id": 分类ID,
      "name": "分类名称"
    }
  ]
  ```

### 上传资源文件
- **方法**: POST
- **路径**: /resources/upload
- **参数**:
  - file (表单文件): 上传的文件
- **响应**:
  ```json
  {
    "success": true,
    "message": "文件上传成功",
    "data": {
      "url": "文件URL"
    }
  }
  ```

### 搜索资源
- **方法**: GET
- **路径**: /resources/search
- **参数**:
  - q: 搜索关键词
  - page (可选): 页码，默认1
  - pageSize (可选): 每页数量，默认10
  - category (可选): 分类ID
- **响应**:
  ```json
  {
    "resources": [资源数组],
    "total": 总数,
    "page": 当前页,
    "pageSize": 每页数量
  }
  ```

### 创建资源
- **方法**: POST
- **路径**: /resources
- **请求体**:
  ```json
  {
    "title": "资源标题",
    "description": "资源描述",
    "category_id": 分类ID,
    "file_path": "文件路径",
    "file_size": 文件大小,
    "file_type": "文件类型",
    "points_required": 所需积分
  }
  ```
- **响应**:
  ```json
  {
    "id": 资源ID,
    "title": "资源标题",
    "description": "资源描述",
    "status": "pending",
    "created_at": "创建时间"
  }
  ```

### 更新资源
- **方法**: PUT
- **路径**: /resources/:id
- **请求体**:
  ```json
  {
    "title": "新标题",
    "description": "新描述",
    "category_id": 新分类ID,
    "points_required": 新所需积分
  }
  ```
- **响应**:
  ```json
  {
    "id": 资源ID,
    "title": "更新后的标题",
    "description": "更新后的描述",
    "status": "pending",
    "updated_at": "更新时间"
  }
  ```

### 删除资源
- **方法**: DELETE
- **路径**: /resources/:id
- **响应**:
  ```json
  {
    "message": "资源已删除"
  }
  ```

### 获取用户资源
- **方法**: GET
- **路径**: /users/resources
- **参数**:
  - page (可选): 页码，默认1
  - pageSize (可选): 每页数量，默认10
  - status (可选): 资源状态
- **响应**:
  ```json
  {
    "resources": [资源数组],
    "total": 总数,
    "page": 当前页,
    "pageSize": 每页数量
  }
  ```

## 管理员控制器 (AdminController)

### 获取待审核资源
- **方法**: GET
- **路径**: /admin/resources/pending
- **参数**:
  - page (可选): 页码，默认1
  - pageSize (可选): 每页数量，默认10
- **响应**:
  ```json
  {
    "resources": [资源数组],
    "total": 总数,
    "page": 当前页,
    "pageSize": 每页数量
  }
  ```

### 审核资源
- **方法**: POST
- **路径**: /admin/resources/:id/review
- **请求体**:
  ```json
  {
    "status": "approved/rejected",
    "message": "审核消息"
  }
  ```
- **响应**:
  ```json
  {
    "id": 资源ID,
    "status": "新状态",
    "updated_at": "更新时间"
  }
  ```

### 获取用户统计
- **方法**: GET
- **路径**: /admin/stats/users
- **响应**:
  ```json
  {
    "user_count": 用户总数,
    "resource_count": 资源总数,
    "pending_count": 待审核资源数,
    "today_user_count": 今日新增用户数
  }
  ```