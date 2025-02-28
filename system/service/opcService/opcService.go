/*
 * @Date: 2022-03-02 10:32:10
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:markdown
 * @LastEditTime: 2023-12-20 11:54:22
 * @FilePath: \opcConnector\system\service\opcService\opcService.go
 */

package opcService

import (
	// "fmt"

	"fmt"
	"opcConnector/system/core/config"
	"opcConnector/system/core/log"

	// "github.com/gin-gonic/gin"
	//"opcConnector/system/model/RequestModel"

	"github.com/spf13/viper"
	"opcConnector/system/util/opc"
	//"sort"
	"context"
	"time"
)

//保持连接
var opcConnClient opc.Connection = nil

type OpcService struct {
}

func (_this *OpcService) ServerList() []string {
	defer func() {
		if err := recover(); err != nil {
			log.Write(log.Error, "opc服务异常！ServerList")
		}

	}()
	obj := opc.NewAutomationObject()
	s := obj.GetOPCServers(config.Instance().Config.App.OpcHost)
	return s
}

//设置服务
func (_this *OpcService) SetServer(r string) error {
	defer func() {
		if err := recover(); err != nil {
			log.Write(log.Error, "opc服务异常！SetServer")
		}

	}()
	c := config.Instance()
	v := viper.New()
	v.SetConfigFile(config.Configpath)
	v.SetConfigType("toml")
	c.Config.App.OpcServer = r
	v.Set("Config", c.Config)
	err := v.WriteConfig()
	return err
}

//获取点位的值
func (_this *OpcService) Read(r []string) (result_err error, result map[string]opc.Item) {
	defer func() {
		if err := recover(); err != nil {
			opcConnClient = nil //清空连接
			result_err = fmt.Errorf("tag异常或服务异常")
			log.Write(log.Error, "opc服务异常！"+fmt.Sprint(err))
		}
	}()

	// 添加超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resultChan := make(chan map[string]opc.Item, 1)
	errChan := make(chan error, 1)

	go func() {
		// 不再使用保持连接模式，每次都创建新连接
		client, err := opc.NewConnection(
			config.Instance().Config.App.OpcServer,
			[]string{config.Instance().Config.App.OpcHost},
			r,
		)
		if err != nil {
			errChan <- err
			return
		}
		defer client.Close()
		result = client.Read()
		resultChan <- result
	}()

	// 等待结果或超时
	select {
	case <-ctx.Done():
		return fmt.Errorf("OPC读取超时"), nil
	case err := <-errChan:
		return err, nil
	case result := <-resultChan:
		return nil, result
	}
}

//写入
func (_this *OpcService) Write(r map[string]interface{}) (result_err error, result_map map[string]bool) {
	result_map = make(map[string]bool)
	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("tag异常或服务异常：%s", fmt.Sprint(err))
			log.Write(log.Error, "opc服务异常！"+fmt.Sprint(err))
		}
	}()

	// 创建新连接
	client, err := opc.NewConnection(
		config.Instance().Config.App.OpcServer,
		[]string{config.Instance().Config.App.OpcHost},
		nil,
	)
	if err != nil {
		return err, nil
	}
	defer client.Close()

	// 写入所有标签
	for k, v := range r {
		client.Add(k)
		err := client.Write(k, v)
		if err != nil {
			result_map[k] = false
			result_err = err
		} else {
			result_map[k] = true
		}
	}
	return result_err, result_map
}

func (_this *OpcService) TagTree() (result_err error, v interface{}) {
	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("tag异常或服务异常")
			log.Write(log.Error, "opc服务异常！"+err.(error).Error())
		}

	}()

	browser, result_err := opc.CreateBrowser(
		config.Instance().Config.App.OpcServer,         // ProgId
		[]string{config.Instance().Config.App.OpcHost}, //  OPC servers nodes
	)
	// PrettyPrint(browser)
	if result_err != nil {
		return result_err, nil
	}
	var opcNameResult OpcName
	result_err, v = _this.TreeToOpcName(browser, opcNameResult)
	return
}

//查询子集
//2023 07 11 很难一个层级一个层级的查询
func (_this *OpcService) GetChildren(path []string) (re []opc.Leaf, result_err error) {
	defer func() {
		if err := recover(); err != nil {
			result_err = fmt.Errorf("tag异常或服务异常")
			log.Write(log.Error, "opc服务异常！"+err.(error).Error())
		}
	}()

	return opc.GetChildNode(config.Instance().Config.App.OpcServer, // ProgId
		[]string{config.Instance().Config.App.OpcHost}, //  OPC servers nodes)
		path)

}

type OpcName struct {
	Name      string    `json:"title"`
	Type      string    `json:"type"`
	DataType  int16     `json:"dataType"`
	Chrildren []OpcName `json:"children"`
}

//递归  *opc.Tree
func (_this *OpcService) TreeToOpcName(tree *opc.Tree, in OpcName) (error, OpcName) {
	in.Name = tree.Name
	in.Type = "dir"

	for _, l := range tree.Leaves {
		in.Chrildren = append(in.Chrildren, OpcName{Name: l.Tag, Type: "tag", DataType: l.Type})
	}

	if tree.Branches != nil {
		for _, t := range tree.Branches {
			var opcNameResult OpcName
			_, re := _this.TreeToOpcName(t, opcNameResult)
			in.Chrildren = append(in.Chrildren, re)
		}
	}

	return nil, in
}

//PrettyPrint prints tree in a nice format
func PrettyPrint(tree *opc.Tree) {
	fmt.Println(tree.Name)
	printSubtree(tree, 1)
}

// printSubtree is a recursive helper function to traverse the tree
func printSubtree(tree *opc.Tree, level int) {
	space := ""
	for i := 0; i < level; i++ {
		space += "  "
	}
	for _, l := range tree.Leaves {
		fmt.Println(space, "-", l.Tag)
	}
	for _, b := range tree.Branches {
		fmt.Println(space, "+", b.Name)
		printSubtree(b, level+1)
	}
}
