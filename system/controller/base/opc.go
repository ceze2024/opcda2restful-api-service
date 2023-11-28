/*
 * @Date: 2022-04-12 09:06:05
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2023-07-12 17:31:34
 * @FilePath: \opcConnector\system\controller\base\opc.go
 */
package base

import (
	"opcConnector/system/core/config"
	"opcConnector/system/core/response"
	"opcConnector/system/service/opcService"

	"github.com/gin-gonic/gin"

	// "opcConnector/system/service/md"
	"opcConnector/system/model/RequestModel"
)

var opcSer opcService.OpcService

func Init(c *gin.Context) {
	re := make(map[string]string)
	re["Secret"] = config.Instance().Config.App.Secret
	re["OpcServer"] = config.Instance().Config.App.OpcServer
	response.OkWithData(re, c)
}

//所有配置
func ServerList(c *gin.Context) {
	response.OkWithData(opcSer.ServerList(), c)
}

func SetServer(c *gin.Context) {
	var r RequestModel.OpcServer
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(opcSer.SetServer(r.Name), c)
}

func Read(c *gin.Context) {
	var r RequestModel.OpcTags
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, v := opcSer.Read(r.Tags)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(v, c)
}

func ReadValue(c *gin.Context) {
	var r RequestModel.OpcTags
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, v := opcSer.Read(r.Tags)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var re map[string]interface{} = make(map[string]interface{})
	for key, item := range v {
		re[key] = item.Value
	}

	response.OkWithData(re, c)
}

func Write(c *gin.Context) {
	var r map[string]interface{}
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = opcSer.Write(r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}

func GetChildren(c *gin.Context) {
	var r RequestModel.OpcTreePath
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	v, err := opcSer.GetChildren(r.Path)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(v, c)
}

func TagTree(c *gin.Context) {
	_, v := opcSer.TagTree()
	response.OkWithData(v, c)
}
