# AI聊天功能说明

## 功能概述

本系统集成了基于火山引擎（VolcEngine）的AI大模型聊天功能，支持用户创建多个聊天会话，并在会话中与AI进行对话交流。系统会保存用户的聊天历史记录，方便用户查看和管理。

## 配置说明

系统使用以下环境变量配置AI大模型：

```
AI_API_KEY=e41deff8-b5e8-4000-a588-ea5171dba541
AI_MODEL_ID=deepseek-r1-distill-qwen-7b-250120
AI_BASE_URL=https://api.volcengine.com/v1/llm
```

这些配置已添加到项目的`.env`文件中，您可以根据需要进行修改。

## API接口说明

所有AI聊天相关的API接口都需要用户登录认证，请在请求头中添加`Authorization: Bearer {token}`。

### 1. 创建聊天会话

- **URL**: `/api/chat/sessions`
- **方法**: POST
- **请求体**:
  ```json
  {
    "title": "会话标题"
  }
  ```
- **响应**:
  ```json
  {
    "id": 1,
    "user_id": 123,
    "title": "会话标题",
    "created_at": "2023-06-01T12:00:00Z",
    "updated_at": "2023-06-01T12:00:00Z"
  }
  ```

### 2. 获取会话列表

- **URL**: `/api/chat/sessions`
- **方法**: GET
- **查询参数**: `page=1&pageSize=10`
- **响应**:
  ```json
  {
    "sessions": [
      {
        "id": 1,
        "user_id": 123,
        "title": "会话标题",
        "created_at": "2023-06-01T12:00:00Z",
        "updated_at": "2023-06-01T12:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 10
  }
  ```

### 3. 获取会话消息历史

- **URL**: `/api/chat/sessions/:id/messages`
- **方法**: GET
- **查询参数**: `page=1&pageSize=50`
- **响应**:
  ```json
  {
    "messages": [
      {
        "id": 1,
        "session_id": 1,
        "user_id": 123,
        "role": "user",
        "content": "你好，AI",
        "created_at": "2023-06-01T12:01:00Z"
      },
      {
        "id": 2,
        "session_id": 1,
        "user_id": 123,
        "role": "assistant",
        "content": "你好！有什么我可以帮助你的吗？",
        "created_at": "2023-06-01T12:01:05Z"
      }
    ],
    "total": 2,
    "page": 1,
    "pageSize": 50
  }
  ```

### 4. 删除会话

- **URL**: `/api/chat/sessions/:id`
- **方法**: DELETE
- **响应**:
  ```json
  {
    "message": "会话已删除"
  }
  ```

### 5. 发送消息

- **URL**: `/api/chat/messages`
- **方法**: POST
- **请求体**:
  ```json
  {
    "session_id": 1,
    "content": "你好，AI"
  }
  ```
- **响应**:
  ```json
  {
    "id": 2,
    "session_id": 1,
    "user_id": 123,
    "role": "assistant",
    "content": "你好！有什么我可以帮助你的吗？",
    "created_at": "2023-06-01T12:01:05Z"
  }
  ```

## 数据模型

系统使用以下数据模型存储聊天相关信息：

### ChatSession（聊天会话）

```go
type ChatSession struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Messages  []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:SessionID"`
}
```

### ChatMessage（聊天消息）

```go
type ChatMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SessionID uint      `json:"session_id" gorm:"index"`
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role" gorm:"type:varchar(20)"` // user 或 assistant
	Content   string    `json:"content" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}
```