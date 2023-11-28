/*
 * @Date: 2022-02-14 12:00:49
 * @LastEditors: 春贰
 * @Desc:
 * @LastEditTime: 2023-07-11 13:53:52
 * @FilePath: \opcConnector\system\core\config\config.go
 */
package config

import (
	"fmt"
	"io/ioutil" //
	"log"

	"opcConnector/system/common/initial"
	"opcConnector/system/util/file"
	"opcConnector/system/util/sys"

	"github.com/spf13/viper" //配置
)

var c *conf

var Configpath string

func Instance() *conf {

	basedir := sys.ExecutePath() + "\\" //根目录
	if c == nil {
		InitConfig(basedir + "config.toml")
	}
	return c
}

type conf struct {
	Config Config
}

type Config struct {
	App    AppConf
	Zaplog ZapLogConf
}

type AppConf struct {
	Secret    string
	HttpPort  int `json:"http-port"`
	OpcHost   string
	OpcServer string
	KeepConn  int
}

type ZapLogConf struct {
	Director string ` json:"director"`
	Level    string ` json:"level" `
}

func InitConfig(tomlPath ...string) *conf {
	if len(tomlPath) > 1 {
		log.Fatal("配置路径数量不正确")
	}
	if file.CheckNotExist(tomlPath[0]) {
		err := ioutil.WriteFile(tomlPath[0], []byte(initial.ConfigToml), 0777)
		if err != nil {
			log.Fatal("无法写入配置模板", err.Error())
		}
		log.Fatal("首次启动！生成【config.toml】，请按需修改后重新启动程序！")
	}
	v := viper.New()
	Configpath = tomlPath[0] //初始化Configpath位置
	v.SetConfigFile(Configpath)
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("配置文件读取失败: ", err.Error())
	}
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("配置解析失败:", err.Error())
	}
	return c
}

//设置配置文件
func SetConfigFile() {

	v := viper.New()
	v.SetConfigFile(Configpath)
	v.SetConfigType("toml")
	c.Config.App.HttpPort = 9022
	v.Set("Config", c.Config)
	err := v.WriteConfig()
	fmt.Println(err)
}
