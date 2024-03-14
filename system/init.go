/*
 * @Date: 2022-08-01 14:47:52
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:启动核心服务
 * @LastEditTime: 2024-03-14 16:46:47
 * @FilePath: \opcConnector\system\init.go
 */
package system

import (
	"opcConnector/system/core/config"
	"opcConnector/system/core/log"
	"opcConnector/system/core/task"
	"opcConnector/system/router"
	"opcConnector/system/util/opc"
	"opcConnector/system/util/sys"
	"os"

	"context"
	"embed"
	"fmt"
	"net/http"

	//"os/exec"
	"os/exec"
	"os/signal"
	"syscall"
)

var basedir = sys.ExecutePath() + "/" //根目录

func Init(staticFs embed.FS) {

	config.Instance() //初始化配置文件

	log.InitLog() //初始化日志

	temlog := ""
	for _, v := range os.Args {
		temlog += v + " "
	}
	log.Write(log.Info, temlog)
	opcTestInit() //验证opc是否注册ok

	r := router.InitRouter(staticFs)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Instance().Config.App.HttpPort),
		Handler:      r,
		ReadTimeout:  0, //设置超时时间
		WriteTimeout: 0, //设置超时时间
		//1GB
		MaxHeaderBytes: 1 << 30,
	}

	url := fmt.Sprintf(`http://127.0.0.1:%d`, config.Instance().Config.App.HttpPort)
	fmt.Println("运行地址：" + url)
	exec.Command("cmd", "/C", "start "+url).Run() //windows打开默认浏览器
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Write(log.Error, err.Error())
			os.Exit(0)
		}
	}()
	go task.Init() //初始化system任务服务

	shutDown(s)
}

func shutDown(s *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Write(log.Fatal, "服务关闭:"+err.Error())
	}
	fmt.Println("退出服务")
}

//验证opc是否注册ok
func opcTestInit() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("opc服务异常！", err)
			log.Write(log.Error, "opc服务异常！")
		}

	}()

	//exec.Command("cmd", "/C", "regsvr32 gbda_aut.dll").Run() //注册dll

	obj := opc.NewAutomationObject()
	obj.GetOPCServers("localhost")
	defer obj.Close()

}
