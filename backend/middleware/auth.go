package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"g/front/backend/config"
	"g/front/backend/models"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式无效"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
			}

			// 返回密钥
			return []byte(config.GetEnv("JWT_SECRET", "your-secret-key")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 验证token是否有效
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查token是否过期
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌已过期"})
				c.Abort()
				return
			}

			// 将用户ID存储在上下文中
			userID := uint(claims["user_id"].(float64))
			c.Set("userID", userID)

			// 继续处理请求
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(user models.User) (string, error) {
	// 设置过期时间（例如24小时）
	expTime := time.Now().Add(24 * time.Hour)

	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expTime.Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET", "your-secret-key")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
