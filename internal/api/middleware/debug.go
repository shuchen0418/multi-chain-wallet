package middleware

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// DebugRequest 记录请求详情，帮助排查参数问题
func DebugRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对特定路径进行调试
		path := c.Request.URL.Path

		// 如果是导入钱包或创建钱包路径，记录请求详情
		if path == "/api/v1/wallet/import/mnemonic" ||
			path == "/api/v1/wallet/import/privatekey" ||
			path == "/api/v1/wallet/create" {
			// 打印请求Method和Path
			log.Printf("[DEBUG] %s %s\n", c.Request.Method, path)

			// 打印请求头
			log.Println("[DEBUG] Headers:")
			for k, v := range c.Request.Header {
				log.Printf("[DEBUG]   %s: %v\n", k, v)
			}

			// 打印请求体
			if c.Request.Method == "POST" || c.Request.Method == "PUT" {
				// 保存请求体，因为读取后会清空
				var bodyBytes []byte
				if c.Request.Body != nil {
					bodyBytes, _ = io.ReadAll(c.Request.Body)
				}

				// 重新设置请求体，以便后续处理
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// 打印请求体
				log.Printf("[DEBUG] Body: %s\n", string(bodyBytes))
			}

			// 打印URL查询参数
			if len(c.Request.URL.RawQuery) > 0 {
				log.Printf("[DEBUG] Query: %s\n", c.Request.URL.RawQuery)
			}
		}

		c.Next()
	}
}
