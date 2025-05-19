# API 文档

本文档提供了项目所有 API 端点的详细信息。

## 认证

需要认证的接口，请求头中需要携带 `Authorization` 字段，值为 `Bearer <token>`。

## 模块

- [用户模块](#用户模块)
- [资源模块](#资源模块)
- [论坛模块](#论坛模块)
- [聊天模块](#聊天模块)
- [积分模块](#积分模块)
- [管理模块](#管理模块)

## 用户模块

### 1. 用户注册

- **描述**: 创建一个新用户账户。
- **方法**: `POST`
- **路径**: `/api/register`
- **认证**: 否
- **请求体 (JSON)**:
  ```json
  {
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }
  ```
- **请求参数说明**:
  - `username` (string, required): 用户名，必须唯一。
  - `email` (string, required): 电子邮箱，必须唯一且格式正确。
  - `password` (string, required): 密码，最小长度6位。
- **成功响应 (201 Created)**:
  ```json
  {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "points": 0,
      "role": "user"
    },
    "token": "your_jwt_token_here"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误 (例如：用户名已存在，邮箱已注册，密码长度不足等)。
    ```json
    {
      "error": "用户名已被使用"
    }
    ```
  - `500 Internal Server Error`: 服务器内部错误 (例如：密码加密失败，创建用户失败，生成令牌失败)。
    ```json
    {
      "error": "创建用户失败"
    }
    ```

### 2. 用户登录

- **描述**: 用户使用用户名/邮箱和密码登录。
- **方法**: `POST`
- **路径**: `/api/login`
- **认证**: 否
- **请求体 (JSON)**:
  ```json
  {
    "username": "testuser", // 或者 test@example.com
    "password": "password123"
  }
  ```
- **请求参数说明**:
  - `username` (string, required): 用户名或电子邮箱。
  - `password` (string, required): 密码。
- **成功响应 (200 OK)**:
  ```json
  {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "avatar": "url_to_avatar_or_null",
      "points": 100,
      "role": "user"
    },
    "token": "your_jwt_token_here"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 用户名或密码错误。
    ```json
    {
      "error": "用户名或密码错误"
    }
    ```
  - `500 Internal Server Error`: 生成令牌失败。
    ```json
    {
      "error": "生成令牌失败"
    }
    ```

### 3. 获取用户资料

- **描述**: 获取当前登录用户的个人资料。
- **方法**: `GET`
- **路径**: `/api/user/profile`
- **认证**: 是
- **成功响应 (200 OK)**:
  ```json
  {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "avatar": "url_to_avatar_or_null",
    "points": 100,
    "role": "user"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 用户不存在。

### 4. 更新用户资料

- **描述**: 更新当前登录用户的个人资料。
- **方法**: `PUT`
- **路径**: `/api/user/profile`
- **认证**: 是
- **请求体 (JSON)**:
  ```json
  {
    "username": "newusername", // 可选
    "email": "newemail@example.com" // 可选
  }
  ```
- **请求参数说明**:
  - `username` (string, optional): 新的用户名，如果提供，必须唯一。
  - `email` (string, optional): 新的电子邮箱，如果提供，必须唯一且格式正确。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "用户资料更新成功",
    "user": {
      "id": 1,
      "username": "newusername",
      "email": "newemail@example.com",
      "avatar": "url_to_avatar_or_null",
      "points": 100,
      "role": "user"
    }
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误 (例如：用户名已存在，邮箱已注册)。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 用户不存在。
  - `500 Internal Server Error`: 更新失败。

### 5. 上传/更新用户头像

- **描述**: 上传或更新当前登录用户的头像。
- **方法**: `POST` 或 `PUT`
- **路径**: `/api/user/avatar`
- **认证**: 是
- **请求体 (form-data)**:
  - `avatar` (file, required): 头像图片文件。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "头像上传成功",
    "avatar_url": "url_to_new_avatar"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 未上传文件或文件类型不支持。
  - `401 Unauthorized`: 未授权访问。
  - `500 Internal Server Error`: 文件上传失败或保存失败。

### 6. 刷新JWT令牌

- **描述**: 使用有效的旧令牌获取新的JWT令牌。
- **方法**: `POST`
- **路径**: `/api/user/refresh-token`
- **认证**: 是 (需要有效的旧令牌)
- **成功响应 (200 OK)**:
  ```json
  {
    "token": "your_new_jwt_token_here"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权或旧令牌无效。
  - `404 Not Found`: 用户不存在。
  - `500 Internal Server Error`: 生成新令牌失败。

## 资源模块

### 1. 获取资源列表

- **描述**: 获取公开的、已审核通过的资源列表，支持分页、分类和排序。
- **方法**: `GET`
- **路径**: `/api/resources`
- **认证**: 否
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
  - `category` (string, optional): 分类ID，用于筛选特定分类的资源。
  - `sort` (string, optional, default: 'newest'): 排序方式。可选值: 'newest' (最新), 'popular' (热门，按下载量)。
  - `query` (string, optional): 搜索关键词，用于按标题或描述搜索资源 (至少2个字符)。
- **成功响应 (200 OK)**:
  ```json
  {
    "resources": [
      {
        "id": 1,
        "title": "示例资源",
        "description": "这是一个示例资源的描述。",
        "category_id": 1,
        "category": {
          "id": 1,
          "name": "技术分享"
        },
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "uploader"
        },
        "file_path": "path/to/resource.zip",
        "file_size": 10240,
        "download_count": 100,
        "status": "approved",
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10,
    "query": "关键词",
    "category": "1",
    "sort": "newest"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 搜索关键词过短。

### 2. 获取资源分类列表

- **描述**: 获取所有可用的资源分类。
- **方法**: `GET`
- **路径**: `/api/resources/categories`
- **认证**: 否
- **成功响应 (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "name": "技术分享",
      "description": "关于各种技术的分享"
    },
    {
      "id": 2,
      "name": "学习资料",
      "description": "各种学习相关的资料"
    }
  ]
  ```

### 3. 获取指定ID的资源详情

- **描述**: 获取特定ID的资源详细信息。
- **方法**: `GET`
- **路径**: `/api/resources/:id`
- **认证**: 否
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "id": 1,
    "title": "示例资源",
    "description": "这是一个示例资源的描述。",
    "category_id": 1,
    "category": {
      "id": 1,
      "name": "技术分享"
    },
    "user_id": 1,
    "user": {
      "id": 1,
      "username": "uploader"
    },
    "file_path": "path/to/resource.zip",
    "file_size": 10240,
    "download_count": 100,
    "status": "approved",
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
  ```
- **错误响应**:
  - `404 Not Found`: 资源不存在。

### 4. 搜索资源

- **描述**: 根据关键词、分类、标签等条件搜索资源。
- **方法**: `GET`
- **路径**: `/api/resources/search`
- **认证**: 否
- **查询参数**:
  - `q` (string, optional): 搜索关键词 (至少2个字符)。
  - `category_id` (integer, optional): 分类ID。
  - `tags` (string, optional): 标签，多个用逗号分隔。
  - `sort` (string, optional, default: 'created_at:desc'): 排序方式。可选值: 'created_at:desc', 'download_count:desc', 'title:asc', 'rating:desc'。
  - `price_range` (string, optional, default: 'all'): 价格范围。可选值: 'all', 'free', 'paid'。
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 12): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "resources": [
      // ... 资源列表 ...
    ],
    "total": 50,
    "page": 1,
    "pageSize": 12,
    "query": "关键词",
    "category": "1",
    "sort": "created_at:desc"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 搜索关键词过短。

### 5. 获取指定资源的评论列表

- **描述**: 获取特定资源下的评论列表，支持分页。
- **方法**: `GET`
- **路径**: `/api/resources/:id/comments`
- **认证**: 否
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "comments": [
      {
        "id": 1,
        "user_id": 2,
        "user": {
          "id": 2,
          "username": "commenter"
        },
        "resource_id": 1,
        "content": "这个资源很棒！",
        "time": "2023-10-28T10:00:00Z"
      }
    ],
    "total": 20,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 资源ID不能为空。

### 6. 删除指定ID的评论 (用户或管理员)

- **描述**: 删除指定ID的评论。评论所有者或管理员可以删除。
- **方法**: `DELETE`
- **路径**: `/api/resources/comments/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 评论ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "评论删除成功"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权。
  - `403 Forbidden`: 无权删除此评论。
  - `404 Not Found`: 评论不存在或用户不存在。
  - `500 Internal Server Error`: 删除评论失败。

### 7. 删除用户自己的指定ID资源

- **描述**: 用户删除自己上传的指定ID的资源。
- **方法**: `DELETE`
- **路径**: `/api/user/resources/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "资源删除成功"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 无效的资源ID。
  - `401 Unauthorized`: 用户未认证。
  - `404 Not Found`: 资源不存在或不属于当前用户。
  - `500 Internal Server Error`: 查询资源失败或删除资源失败。

### 8. 获取用户收藏的资源列表

- **描述**: 获取当前登录用户收藏的所有资源。
- **方法**: `GET`
- **路径**: `/api/user/favorites/resources` (或 `/api/favorites/resources`)
- **认证**: 是
- **成功响应 (200 OK)**:
  ```json
  {
    "data": [
      {
        "id": 1,
        "title": "收藏的资源1",
        "description": "描述...",
        "category": {
          "id": 1,
          "name": "技术"
        },
        "user": {
          "id": 5,
          "username": "uploader_user"
        },
        "created_at": "2023-10-28T12:00:00Z" // 收藏时间
      }
    ]
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `500 Internal Server Error`: 获取收藏列表失败。

### 9. 发表对指定资源的评论

- **描述**: 用户对指定的资源发表评论。
- **方法**: `POST`
- **路径**: `/api/resources/:id/comments`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **请求体 (JSON)**:
  ```json
  {
    "content": "这是一条评论内容。"
  }
  ```
- **请求参数说明**:
  - `content` (string, required): 评论内容。
- **成功响应 (201 Created)**:
  ```json
  {
    "id": 10,
    "user_id": 1,
    "resource_id": 1,
    "content": "这是一条评论内容。",
    "time": "2023-10-28T14:00:00Z",
    "user": {
      "id": 1,
      "username": "current_user"
    }
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 资源ID不能为空或请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 创建评论失败。

### 10. 创建新资源 (上传资源信息)

- **描述**: 用户创建新的资源条目 (不包含文件上传，文件上传有单独接口)。
- **方法**: `POST`
- **路径**: `/api/resources`
- **认证**: 是
- **请求体 (JSON)**:
  ```json
  {
    "title": "我的新资源",
    "description": "这是新资源的详细描述。",
    "category_id": 1,
    "tags": ["tag1", "tag2"], // 可选
    "price": 0 // 可选，默认为0
  }
  ```
- **请求参数说明**:
  - `title` (string, required): 资源标题。
  - `description` (string, required): 资源描述。
  - `category_id` (integer, required): 资源分类ID。
  - `tags` (array of strings, optional): 资源标签。
  - `price` (float, optional): 资源价格。
- **成功响应 (201 Created)**:
  ```json
  {
    "id": 15,
    "title": "我的新资源",
    // ... 其他资源信息 ...
    "status": "pending" // 新资源默认为待审核
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 分类不存在。
  - `500 Internal Server Error`: 创建资源失败。

### 11. 更新指定ID的资源信息

- **描述**: 用户更新自己上传的资源的信息 (通常是管理员或资源所有者)。
- **方法**: `PUT`
- **路径**: `/api/resources/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **请求体 (JSON)**:
  ```json
  {
    "title": "更新后的资源标题", // 可选
    "description": "更新后的描述。", // 可选
    "category_id": 2, // 可选
    "tags": ["newtag"], // 可选
    "price": 10 // 可选
  }
  ```
- **成功响应 (200 OK)**:
  ```json
  {
    "id": 15,
    "title": "更新后的资源标题",
    // ... 其他更新后的资源信息 ...
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权修改此资源。
  - `404 Not Found`: 资源或分类不存在。
  - `500 Internal Server Error`: 更新资源失败。

### 12. 删除指定ID的资源 (通常是管理员或资源所有者)

- **描述**: 删除指定ID的资源。
- **方法**: `DELETE`
- **路径**: `/api/resources/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "资源删除成功"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权删除此资源。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 删除资源失败。

### 13. 获取用户自己上传的资源列表

- **描述**: 获取当前登录用户上传的所有资源 (包括各种状态的资源)。
- **方法**: `GET`
- **路径**: `/api/user/resources` (或 `/api/user/my-resources`)
- **认证**: 是
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
  - `status` (string, optional): 资源状态 (e.g., 'approved', 'pending', 'rejected')。
- **成功响应 (200 OK)**:
  ```json
  {
    "resources": [
      // ... 用户上传的资源列表 ...
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `500 Internal Server Error`: 获取资源列表失败。

### 14. 用户删除自己上传的指定ID的资源 (特定路径)

- **描述**: 用户删除自己上传的指定ID的资源 (通过 `/api/user/my-resources/:id` 路径)。
- **方法**: `DELETE`
- **路径**: `/api/user/my-resources/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "资源删除成功"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 无效的资源ID。
  - `401 Unauthorized`: 用户未认证。
  - `404 Not Found`: 资源不存在或不属于当前用户。
  - `500 Internal Server Error`: 删除资源失败。

### 15. 添加资源到收藏夹

- **描述**: 用户将指定ID的资源添加到自己的收藏夹。
- **方法**: `POST`
- **路径**: `/api/resources/:id/favorite`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (201 Created)**:
  ```json
  {
    "success": true,
    "isFavorited": true,
    "message": "收藏成功"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 已收藏该资源。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 收藏失败。

### 16. 从收藏夹移除资源

- **描述**: 用户将指定ID的资源从自己的收藏夹中移除。
- **方法**: `DELETE`
- **路径**: `/api/resources/:id/favorite`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "已取消收藏",
    "isFavorited": false
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 取消收藏失败。

### 17. 获取指定资源的收藏状态

- **描述**: 用户查询自己是否收藏了指定ID的资源。
- **方法**: `GET`
- **路径**: `/api/resources/:id/favorite-status`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "is_favorite": true, // 或者 false
    "isFavorited": true // 或者 false (建议统一字段名)
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。

### 18. 上传资源文件

- **描述**: 用户上传资源文件。通常在创建资源信息后调用此接口关联文件。
- **方法**: `POST`
- **路径**: `/api/resources/upload` (或者可以设计为 `/api/resources/:id/upload` 来关联特定资源)
- **认证**: 是
- **请求体 (form-data)**:
  - `file` (file, required): 资源文件。
  - `resource_id` (integer, optional): 如果路径不包含:id，则需要此参数来关联资源。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "文件上传成功",
    "file_path": "uploaded/path/to/file.zip",
    "file_url": "http://minio-server/bucket/uploaded/path/to/file.zip"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 未上传文件或文件类型不支持或resource_id无效。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 如果需要resource_id但未找到对应资源。
  - `500 Internal Server Error`: 文件上传或保存失败。

### 19. 获取资源下载链接

- **描述**: 获取指定ID资源的实际下载链接 (通常是MinIO的预签名URL或直接链接)。
- **方法**: `GET`
- **路径**: `/api/download/:id`
- **认证**: 是 (或者根据业务逻辑决定是否需要认证，例如付费资源)
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "url": "http://minio-server/bucket/path/to/resource.zip?signature=...",
    "filename": "resource.zip"
  }
  ```
- **错误响应**:
  - `404 Not Found`: 资源不存在或文件在存储中不存在。

### 20. 点赞资源

- **描述**: 用户点赞指定ID的资源。
- **方法**: `POST`
- **路径**: `/api/resources/:id/like`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "点赞成功",
    "liked": true,
    "likes": 101,
    "dislikes": 5
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 已点赞或已点踩 (如果业务逻辑是互斥的)。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 点赞失败。

### 21. 取消点赞/点踩资源 (点踩)

- **描述**: 用户取消对指定ID资源的赞，或者点踩该资源 (根据实现，这里是dislike)。
- **方法**: `DELETE` (语义上更像是取消点赞，但路径是 /dislike，所以理解为点踩)
- **路径**: `/api/resources/:id/dislike`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "操作成功", // "点踩成功" 或 "取消点踩成功"
    "disliked": true, // 或 false
    "likes": 100,
    "dislikes": 6 // 或 5
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 未点赞/点踩 (如果要取消的操作对象不存在)。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 操作失败。

### 22. 获取指定资源的当前用户点赞状态

- **描述**: 查询当前用户是否点赞/点踩了指定ID的资源。
- **方法**: `GET`
- **路径**: `/api/resources/:id/like-status`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 资源ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "liked": true, // boolean
    "disliked": false // boolean
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 查询状态失败。

## 论坛模块

### 1. 获取论坛分类列表

- **描述**: 获取所有论坛的分类信息，包含各分类下的主题数和帖子数。
- **方法**: `GET`
- **路径**: `/api/forum/categories`
- **认证**: 否
- **成功响应 (200 OK)**:
  ```json
  [
    {
      "id": 1,
      "name": "技术交流",
      "description": "讨论技术的板块",
      "parent_id": null,
      "topic_count": 50,
      "post_count": 200
    }
  ]
  ```

### 2. 获取主题列表

- **描述**: 获取论坛主题列表，支持分页和按分类筛选。
- **方法**: `GET`
- **路径**: `/api/forum/topics`
- **认证**: 否
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
  - `category` (string, optional): 分类ID，用于筛选特定分类下的主题。
- **成功响应 (200 OK)**:
  ```json
  {
    "topics": [
      {
        "id": 1,
        "user_id": 1,
        "user": {
          "id": 1,
          "username": "topic_starter"
        },
        "category_id": 1,
        "category": {
          "id": 1,
          "name": "技术交流"
        },
        "title": "这是一个论坛主题",
        "content": "主题内容...",
        "view_count": 100,
        "reply_count": 10,
        "like_count": 20,
        "dislike_count": 2,
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "pageSize": 10
  }
  ```

### 3. 获取指定ID的主题详情

- **描述**: 获取特定ID的论坛主题的详细信息及其回复列表。
- **方法**: `GET`
- **路径**: `/api/forum/topics/:id`
- **认证**: 否
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **查询参数 (用于回复列表分页)**:
  - `replyPage` (integer, optional, default: 1): 回复列表的页码。
  - `replyPageSize` (integer, optional, default: 10): 回复列表每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "topic": {
      "id": 1,
      "user_id": 1,
      "user": { "id": 1, "username": "starter" },
      "category_id": 1,
      "category": { "id": 1, "name": "技术" },
      "title": "主题标题",
      "content": "主题内容",
      // ...其他主题字段...
      "created_at": "2023-10-28T00:00:00Z"
    },
    "replies": [
      {
        "id": 101,
        "user_id": 2,
        "user": { "id": 2, "username": "replier" },
        "topic_id": 1,
        "content": "回复内容",
        "created_at": "2023-10-28T01:00:00Z"
      }
    ],
    "replyTotal": 20,
    "replyPage": 1,
    "replyPageSize": 10
  }
  ```
- **错误响应**:
  - `404 Not Found`: 主题不存在。

### 4. 获取指定主题的点赞用户列表 (或点赞数)

- **描述**: 获取点赞了特定主题的用户列表或仅点赞数量。
- **方法**: `GET`
- **路径**: `/api/forum/topics/:id/likes`
- **认证**: 否
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "likes": 25 // 或者一个用户对象数组
  }
  ```
- **错误响应**:
  - `404 Not Found`: 主题不存在。

### 5. 获取用户收藏的帖子(主题)列表

- **描述**: 获取当前登录用户收藏的所有论坛主题。
- **方法**: `GET`
- **路径**: `/api/user/favorites/posts`
- **认证**: 是
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "favorites": [
      {
        "id": 1, // 收藏记录ID
        "user_id": 1,
        "topic_id": 10,
        "created_at": "2023-10-28T10:00:00Z",
        "topic": {
          "id": 10,
          "title": "我收藏的主题",
          // ... 其他主题信息
        }
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `500 Internal Server Error`: 获取收藏列表失败。

### 6. 创建新的论坛主题

- **描述**: 用户创建一个新的论坛主题。
- **方法**: `POST`
- **路径**: `/api/forum/topics`
- **认证**: 是
- **请求体 (JSON)**:
  ```json
  {
    "title": "我的新主题",
    "content": "这是主题的详细内容。",
    "category_id": 1
  }
  ```
- **请求参数说明**:
  - `title` (string, required): 主题标题。
  - `content` (string, required): 主题内容。
  - `category_id` (integer, required): 主题所属分类ID。
- **成功响应 (201 Created)**:
  ```json
  {
    "id": 101,
    "user_id": 1,
    "category_id": 1,
    "title": "我的新主题",
    "content": "这是主题的详细内容。",
    // ... 其他主题字段 ...
    "created_at": "2023-10-28T15:00:00Z"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 分类不存在。
  - `500 Internal Server Error`: 创建主题失败。

### 7. 更新指定ID的论坛主题

- **描述**: 用户更新自己发布的论坛主题 (或管理员更新任意主题)。
- **方法**: `PUT`
- **路径**: `/api/forum/topics/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **请求体 (JSON)**:
  ```json
  {
    "title": "更新后的主题标题", // 可选
    "content": "更新后的主题内容。", // 可选
    "category_id": 2 // 可选
  }
  ```
- **成功响应 (200 OK)**:
  ```json
  {
    "id": 101,
    "title": "更新后的主题标题",
    // ... 其他更新后的主题信息 ...
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权修改此主题。
  - `404 Not Found`: 主题或分类不存在。
  - `500 Internal Server Error`: 更新主题失败。

### 8. 删除指定ID的论坛主题

- **描述**: 用户删除自己发布的论坛主题 (或管理员删除任意主题)。
- **方法**: `DELETE`
- **路径**: `/api/forum/topics/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "主题删除成功"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权删除此主题。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 删除主题失败。

### 9. 创建对指定主题的回复

- **描述**: 用户对指定的论坛主题发表回复。
- **方法**: `POST`
- **路径**: `/api/forum/topics/:id/replies`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **请求体 (JSON)**:
  ```json
  {
    "content": "这是一条回复内容。"
  }
  ```
- **请求参数说明**:
  - `content` (string, required): 回复内容。
- **成功响应 (201 Created)**:
  ```json
  {
    "id": 205,
    "user_id": 1,
    "topic_id": 101,
    "content": "这是一条回复内容。",
    "created_at": "2023-10-28T16:00:00Z",
    "user": {
      "id": 1,
      "username": "current_user"
    }
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 创建回复失败。

### 10. 更新指定ID的回复

- **描述**: 用户更新自己发表的回复 (或管理员更新任意回复)。
- **方法**: `PUT`
- **路径**: `/api/forum/replies/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 回复ID。
- **请求体 (JSON)**:
  ```json
  {
    "content": "更新后的回复内容。"
  }
  ```
- **成功响应 (200 OK)**:
  ```json
  {
    "id": 205,
    "content": "更新后的回复内容。",
    // ... 其他更新后的回复信息 ...
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权修改此回复。
  - `404 Not Found`: 回复不存在。
  - `500 Internal Server Error`: 更新回复失败。

### 11. 删除指定ID的回复

- **描述**: 用户删除自己发表的回复 (或管理员删除任意回复)。
- **方法**: `DELETE`
- **路径**: `/api/forum/replies/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 回复ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "回复删除成功"
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权删除此回复。
  - `404 Not Found`: 回复不存在。
  - `500 Internal Server Error`: 删除回复失败。

### 12. 点赞指定ID的论坛主题

- **描述**: 用户点赞一个论坛主题。
- **方法**: `POST`
- **路径**: `/api/forum/topics/:id/like`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "点赞成功",
    "liked": true,
    "likes": 21, // 当前总点赞数
    "dislikes": 2 // 当前总点踩数
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 操作失败。

### 13. 取消点赞指定ID的论坛主题

- **描述**: 用户取消对一个论坛主题的点赞。
- **方法**: `DELETE`
- **路径**: `/api/forum/topics/:id/like`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "取消点赞成功",
    "liked": false,
    "likes": 20,
    "dislikes": 2
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 用户未点赞该主题。
  - `401 Unauthorized`: 未授权。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 操作失败。

### 14. 点踩指定ID的论坛主题

- **描述**: 用户点踩一个论坛主题。
- **方法**: `POST`
- **路径**: `/api/forum/topics/:id/dislike`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "点踩成功",
    "disliked": true,
    "liked": false, // 如果点赞和点踩互斥，点赞状态会改变
    "likes": 20,
    "dislikes": 3
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 操作失败。

### 15. 取消点踩指定ID的论坛主题

- **描述**: 用户取消对一个论坛主题的点踩。
- **方法**: `DELETE`
- **路径**: `/api/forum/topics/:id/dislike`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "取消点踩成功",
    "disliked": false,
    "likes": 20,
    "dislikes": 2
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 用户未点踩该主题。
  - `401 Unauthorized`: 未授权。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 操作失败。

### 16. 添加论坛主题到收藏夹

- **描述**: 用户将指定ID的论坛主题添加到自己的收藏夹。
- **方法**: `POST`
- **路径**: `/api/forum/topics/:id/favorite`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (201 Created)**:
  ```json
  {
    "success": true,
    "isFavorited": true,
    "message": "收藏成功"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 已收藏该主题。
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 收藏失败。

### 17. 从收藏夹移除论坛主题

- **描述**: 用户将指定ID的论坛主题从自己的收藏夹中移除。
- **方法**: `DELETE`
- **路径**: `/api/forum/topics/:id/favorite`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "已取消收藏",
    "isFavorited": false
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 主题不存在。
  - `500 Internal Server Error`: 取消收藏失败。

### 18. 获取指定论坛主题的收藏状态

- **描述**: 用户查询自己是否收藏了指定ID的论坛主题。
- **方法**: `GET`
- **路径**: `/api/forum/topics/:id/favorite-status`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 主题ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "isFavorited": true // 或者 false
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 主题不存在。

## 聊天模块

### 1. 创建新的聊天会话

- **描述**: 用户创建一个新的AI聊天会话。
- **方法**: `POST`
- **路径**: `/api/chat/sessions`
- **认证**: 是
- **请求体 (JSON)**:
  ```json
  {
    "title": "我的聊天会话标题"
  }
  ```
- **请求参数说明**:
  - `title` (string, required): 会话标题。
- **成功响应 (201 Created)**:
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "我的聊天会话标题",
    "created_at": "2023-10-28T10:00:00Z",
    "updated_at": "2023-10-28T10:00:00Z"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `500 Internal Server Error`: 创建会话失败。

### 2. 获取用户的所有聊天会话列表

- **描述**: 获取当前登录用户的所有聊天会话，支持分页。
- **方法**: `GET`
- **路径**: `/api/chat/sessions`
- **认证**: 是
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "sessions": [
      {
        "id": 1,
        "user_id": 1,
        "title": "会话1",
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-28T11:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。

### 3. 获取指定会话的消息历史

- **描述**: 获取特定聊天会话中的消息历史记录，支持分页。
- **方法**: `GET`
- **路径**: `/api/chat/sessions/:id/messages`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 会话ID。
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 50): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "messages": [
      {
        "id": 101,
        "session_id": 1,
        "sender": "user", // 'user' 或 'ai'
        "content": "你好AI",
        "created_at": "2023-10-28T11:00:00Z"
      },
      {
        "id": 102,
        "session_id": 1,
        "sender": "ai",
        "content": "你好！有什么可以帮助你的吗？",
        "created_at": "2023-10-28T11:00:05Z"
      }
    ],
    "total": 20,
    "page": 1,
    "pageSize": 50
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 无效的会话ID。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权访问此会话。
  - `404 Not Found`: 会话不存在。

### 4. 删除指定的聊天会话

- **描述**: 用户删除自己的某个聊天会话及其所有消息。
- **方法**: `DELETE`
- **路径**: `/api/chat/sessions/:id`
- **认证**: 是
- **路径参数**:
  - `id` (integer, required): 会话ID。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "会话已删除"
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 无效的会话ID。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权删除此会话。
  - `404 Not Found`: 会话不存在。

### 5. 发送消息并获取AI回复

- **描述**: 用户在指定的聊天会话中发送一条消息，并获取AI的回复。
- **方法**: `POST`
- **路径**: `/api/chat/messages`
- **认证**: 是
- **请求体 (JSON)**:
  ```json
  {
    "session_id": 1,
    "content": "请问今天天气怎么样？"
  }
  ```
- **请求参数说明**:
  - `session_id` (integer, required): 消息所属的会话ID。
  - `content` (string, required): 用户发送的消息内容。
- **成功响应 (200 OK)**:
  ```json
  {
    "user_message": {
      "id": 103,
      "session_id": 1,
      "sender": "user",
      "content": "请问今天天气怎么样？",
      "created_at": "2023-10-28T12:00:00Z"
    },
    "ai_message": {
      "id": 104,
      "session_id": 1,
      "sender": "ai",
      "content": "很抱歉，我暂时无法查询实时天气信息。",
      "created_at": "2023-10-28T12:00:05Z"
    }
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权访问。
  - `403 Forbidden`: 无权访问此会话。
  - `404 Not Found`: 会话不存在。
  - `500 Internal Server Error`: 发送消息失败或AI服务调用失败。

## 积分模块

### 1. 获取用户当前积分

- **描述**: 获取当前登录用户的总积分数。
- **方法**: `GET`
- **路径**: `/api/user/points`
- **认证**: 是
- **成功响应 (200 OK)**:
  ```json
  {
    "points": 500
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。
  - `404 Not Found`: 用户不存在。

### 2. 获取用户积分历史记录

- **描述**: 获取当前登录用户的积分变动历史记录，支持分页。
- **方法**: `GET`
- **路径**: `/api/user/points/history`
- **认证**: 是
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "records": [
      {
        "id": 1,
        "user_id": 1,
        "points": 20,
        "type": "approve", // 积分变动类型，如 'approve', 'admin', 'download'
        "description": "资源审核通过奖励",
        "resource_id": 5, // 可选，关联的资源ID
        "resource": { // 可选，关联的资源信息
            "id": 5,
            "title": "被审核的资源"
        },
        "created_at": "2023-10-28T09:00:00Z"
      },
      {
        "id": 2,
        "user_id": 1,
        "points": -10,
        "type": "download",
        "description": "下载付费资源",
        "resource_id": 10,
        "created_at": "2023-10-27T15:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权访问。

### 3. 管理员为用户添加积分

- **描述**: 管理员为指定用户手动添加或扣除积分。
- **方法**: `POST`
- **路径**: `/api/admin/points/add` (此路径在adminController中，但功能上属于积分管理)
- **认证**: 是 (管理员权限)
- **请求体 (JSON)**:
  ```json
  {
    "user_id": 2,
    "points": 100, // 正数为增加，负数为扣除
    "description": "活动奖励"
  }
  ```
- **请求参数说明**:
  - `user_id` (integer, required): 目标用户ID。
  - `points` (integer, required): 变动的积分数量。
  - `description` (string, required): 积分变动原因描述。
- **成功响应 (200 OK)**:
  ```json
  {
    "message": "积分添加成功",
    "pointRecord": {
      "id": 3,
      "user_id": 2,
      "points": 100,
      "type": "admin",
      "description": "活动奖励",
      "created_at": "2023-10-28T14:00:00Z"
    },
    "newPoints": 600 // 用户新的总积分
  }
  ```
- **错误响应**:
  - `400 Bad Request`: 请求参数错误。
  - `401 Unauthorized`: 未授权。
  - `403 Forbidden`: 权限不足 (非管理员)。
  - `404 Not Found`: 目标用户不存在。
  - `500 Internal Server Error`: 操作失败。

## 管理模块

所有管理模块的接口都需要管理员权限，并且路径以 `/api/admin` 开头。
认证中间件后会再经过一个管理员权限检查中间件。

### 1. 获取待审核资源列表

- **描述**: 管理员获取等待审核的资源列表，支持分页。
- **方法**: `GET`
- **路径**: `/api/admin/resources/pending`
- **认证**: 是 (管理员权限)
- **查询参数**:
  - `page` (integer, optional, default: 1): 页码。
  - `pageSize` (integer, optional, default: 10): 每页数量。
- **成功响应 (200 OK)**:
  ```json
  {
    "resources": [
      {
        "id": 20,
        "title": "待审核的资源",
        "status": "pending",
        "user_id": 3,
        "user": { "id": 3, "username": "uploader3" },
        "category_id": 1,
        "category": { "id": 1, "name": "技术" },
        "created_at": "2023-10-28T10:00:00Z"
        // ... 其他资源字段
      }
    ],
    "total": 5,
    "page": 1,
    "pageSize": 10
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权。
  - `403 Forbidden`: 权限不足。

### 2. 审核资源

- **描述**: 管理员审核指定的资源，可以批准或拒绝。
- **方法**: `PUT`
- **路径**: `/api/admin/resources/:id/review`
- **认证**: 是 (管理员权限)
- **路径参数**:
  - `id` (integer, required): 待审核的资源ID。
- **请求体 (JSON)**:
  ```json
  {
    "status": "approved", // 'approved' 或 'rejected'
    "message": "审核通过，资源质量很好。" // 可选，审核意见
  }
  ```
- **请求参数说明**:
  - `status` (string, required): 审核结果，必须是 'approved' 或 'rejected'。
  - `message` (string, optional): 审核相关的说明信息。
- **成功响应 (200 OK)** (返回更新后的资源信息):
  ```json
  {
    "id": 20,
    "title": "待审核的资源",
    "status": "approved", // 更新后的状态
    // ... 其他资源字段
  }
  ```
  如果审核通过 (`status: "approved"`)，上传该资源的用户会自动获得20积分奖励。
- **错误响应**:
  - `400 Bad Request`: 请求参数错误 (例如status值无效)。
  - `401 Unauthorized`: 未授权。
  - `403 Forbidden`: 权限不足。
  - `404 Not Found`: 资源不存在。
  - `500 Internal Server Error`: 更新资源状态失败。

### 3. 获取用户统计信息

- **描述**: 管理员获取平台的用户相关统计数据。
- **方法**: `GET`
- **路径**: `/api/admin/stats`
- **认证**: 是 (管理员权限)
- **成功响应 (200 OK)**:
  ```json
  {
    "user_count": 1500,       // 总用户数
    "resource_count": 3000,   // 总资源数
    "pending_count": 15,      // 待审核资源数
    "today_user_count": 25    // 今日新增用户数
  }
  ```
- **错误响应**:
  - `401 Unauthorized`: 未授权。
  - `403 Forbidden`: 权限不足。