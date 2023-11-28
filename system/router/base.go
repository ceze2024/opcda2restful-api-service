/*
 * @Date: 2022-04-12 09:17:10
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2023-07-12 17:29:45
 * @FilePath: \opcConnector\system\router\base.go
 */

package router

import (
	"github.com/gin-gonic/gin"
	"opcConnector/system/controller/base"
	"opcConnector/system/middleware"
)

func BaseRouter(r *gin.Engine) {

	r.POST("/init", base.Init)
	r.Use(middleware.JwtMid())
	{
		r.POST("/ServerList", base.ServerList)
		r.POST("/SetServer", base.SetServer)
		r.POST("/Read", base.Read)
		r.POST("/Write", base.Write)
		r.POST("/TagTree", base.TagTree)
		r.POST("/GetChildren", base.GetChildren)
		//额外
		r.POST("/api/opc/read", base.ReadValue)
		r.POST("/api/opc/write", base.Write)

	}

	// authG := r.Group("tools")
	// //只需要登录
	// authG.Use(middleware.JwtMid())
	// {
	// 	authG.POST("/sequence", tools.Sequence)
	// 	authG.POST("/genpdf", pdf.GenPdf)
	// }

}
