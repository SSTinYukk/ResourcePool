package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// ErrorMiddleware 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer捕获panic
		defer func() {
			if err := recover(); err != nil {
				// 记录错误和堆栈信息
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// 返回500错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal Server Error",
					Code:    http.StatusInternalServerError,
					Message: "服务器内部错误，请稍后再试",
				})

				// 终止请求
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last()

			// 记录错误
			log.Printf("Error: %v", err)

			// 如果响应已经发送，则不再处理
			if c.Writer.Written() {
				return
			}

			// 根据错误类型返回不同的状态码
			switch err.Type {
			case gin.ErrorTypeBind:
				// 请求绑定错误
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error:   "Bad Request",
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				})
			case gin.ErrorTypePrivate:
				// 私有错误，通常是业务逻辑错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal Server Error",
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
			default:
				// 其他错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal Server Error",
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
			}
		}
	}
}

// NotFoundHandler 处理404错误
func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error:   "Not Found",
		Code:    http.StatusNotFound,
		Message: "请求的资源不存在",
	})
}
