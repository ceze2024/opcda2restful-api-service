/*
 * @Date: 2022-03-04 11:31:54
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2022-11-14 10:41:13
 * @FilePath: \opcConnector\system\middleware\cors.go
 */
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 不引入"github.com/gin-gonic/gin"
func CorsMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Length,Content-Type,Authorization,Accept, X-Requested-With,secret")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,secret")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

// 引入"github.com/gin-gonic/gin"
// func CorsMid() gin.HandlerFunc {
// 	config := cors.DefaultConfig()
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
// 	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}
// 	config.AllowOrigins = []string{"*"}
// 	config.AllowCredentials = true
// 	return cors.New(config)
// }
