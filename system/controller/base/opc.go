/*
 * @Date: 2022-04-12 09:06:05
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2024-05-24 12:03:56
 * @FilePath: \opcConnector\system\controller\base\opc.go
 */
package base

import (
	"context"
	"net/http"
	"opcConnector/system/core/config"
	"opcConnector/system/core/response"
	"opcConnector/system/service/opcService"
	"opcConnector/system/util/opc"

	"github.com/gin-gonic/gin"

	// "opcConnector/system/service/md"
	"opcConnector/system/model/RequestModel"
	"time"
)

var opcSer opcService.OpcService

func StartIt(c *gin.Context) {
	c.JSON(http.StatusOK, true)
}

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
	// 添加上下文超时控制
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	
	var r RequestModel.OpcTags
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	// 创建一个完成通道
	done := make(chan struct{})
	var opcErr error
	var v map[string]opc.Item
	
	go func() {
		err, v = opcSer.Read(r.Tags)
		if err != nil {
			opcErr = err
		}
		close(done)
	}()
	
	// 等待OPC读取完成或超时
	select {
	case <-ctx.Done():
		response.FailWithMessage("读取超时", c)
		return
	case <-done:
		if opcErr != nil {
			response.FailWithMessage(opcErr.Error(), c)
			return
		}
		response.OkWithData(v, c)
	}
}

func ReadValue(c *gin.Context) {
	// 类似的超时处理
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	
	var r RequestModel.OpcTags
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	
	done := make(chan struct{})
	var opcErr error
	var v map[string]opc.Item
	
	go func() {
		err, v = opcSer.Read(r.Tags)
		if err != nil {
			opcErr = err
		}
		close(done)
	}()
	
	select {
	case <-ctx.Done():
		response.FailWithMessage("读取超时", c)
		return
	case <-done:
		if opcErr != nil {
			response.FailWithMessage(opcErr.Error(), c)
			return
		}
		
		var re map[string]interface{} = make(map[string]interface{})
		for key, item := range v {
			re[key] = item.Value
		}
		response.OkWithData(re, c)
	}
}

func Write(c *gin.Context) {
	var r map[string]interface{}
	err := c.ShouldBind(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, reMap := opcSer.Write(r)
	if err != nil {
		response.FailWithDetailed(reMap, err.Error(), c)
		return
	}
	response.OkWithData(reMap, c)
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
