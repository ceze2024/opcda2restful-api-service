/*
 * @Date: 2022-02-14 10:46:28
 * @LastEditors: 春贰
 * @Desc:
 * @LastEditTime: 2023-07-12 15:24:39
 * @FilePath: \opcConnector\system\router\InitRouter.go
 */

package router

import (
	"embed"
	"io/fs"
	"net/http"

	"opcConnector/system/middleware"
	"opcConnector/system/util/sys"

	"github.com/gin-gonic/gin"
)

func InitRouter(staticFs embed.FS) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	basedir := sys.ExecutePath() + "\\" //根目录
	r := gin.New()

	r.Static("/runtime/file", basedir+"runtime/file")
	//r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.GET("/", func(c *gin.Context) {
		// c.Request.URL.Path = "/pages/"
		// r.HandleContext(c)
		c.Redirect(http.StatusMovedPermanently, "/pages/")
	})

	// r.LoadHTMLGlob(basedir + "template/*")
	// r.LoadHTMLFiles(basedir + "pages/index.html")
	// r.GET("/opc", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"secret":          config.Instance().Config.App.Secret,
	// 		"activeOpcServer": config.Instance().Config.App.OpcServer,
	// 	})
	// })

	r.Static("/pages1", basedir+"pages") //静态文件

	fads, _ := fs.Sub(staticFs, "pages")
	r.StaticFS("/pages", http.FS(fads)) //挂载到二进制中

	// r.StaticFS("/apidoc", http.Dir(basedir+"apidoc"))

	//r.Use(gin.Logger())
	//r.Use(gin.Recovery())

	//全局中间件
	r.Use(middleware.CorsMid())
	//注册基础路由
	BaseRouter(r)
	return r
}
