package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 额外添加的认证中间件
// 实际应用中应该实现真正的认证逻辑
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取token
		token := c.GetHeader("Authorization")

		// 检查token是否有效
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// 实际应用中应该验证token的有效性
		// 这里仅做示例，假设token总是有效的

		c.Next()
	}
}

// RequestLogger 高级日志中间件，包含更多详细信息
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		fmt.Printf("[GIN] %v | %3d | %13v | %15s | %s | %s\n",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			reqMethod,
			reqURI,
		)
	}
}
