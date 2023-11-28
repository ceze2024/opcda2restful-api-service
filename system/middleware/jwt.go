/*
 * @Date: 2022-02-16 10:07:39
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2023-10-16 16:27:24
 * @FilePath: \opcConnector\system\middleware\jwt.go
 */

package middleware

import (
	"bytes"
	"io/ioutil"
	"opcConnector/system/core/config"
	"opcConnector/system/core/log"
	"opcConnector/system/core/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)


type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

//重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func JwtMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		Secret := c.Request.Header.Get("Secret")

		// data, err := c.GetRawData()     //这玩意读不到
		body, _ := ioutil.ReadAll(c.Request.Body)

		//写入的
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer


		// //把读过的字节流重新放到body
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		
		trueSecret := config.Instance().Config.App.Secret
		//如果没配置Secret，就不校验
		if trueSecret != "" {
			if Secret != trueSecret {
				response.FailWithDetailed(gin.H{"reload": true}, "error Secret", c)
				c.Abort()
				return
			}
		}
		c.Next()

		//获取返回 
		lenbody := len(body)
		responseBody := writer.b.String()
		//请求大于1M数据不记录，避免上传大文件还要处理
		if lenbody < 10<<20 {
			//记录日志
			log.Write(log.Info, "记录请求", zap.String("RequestUrl", c.Request.URL.Path), zap.String("RequestIp", c.ClientIP()), zap.String("RequestUserAgent", c.Request.UserAgent()), zap.String("RequestMethod", c.Request.Method), zap.String("Secret", Secret), zap.String("RequestData", string(body)), zap.String("ReturnData", responseBody))
		}

	
	}
}
