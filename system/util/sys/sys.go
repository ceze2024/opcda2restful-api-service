/*
 * @Date: 2022-08-11 15:50:15
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2023-06-30 15:25:40
 * @FilePath: \opcConnector\system\util\sys\sys.go
 */
package sys

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"opcConnector/system/util/encrypt"
	"opcConnector/system/util/str"
)

func RealPath(f string) string {
	p, err := filepath.Abs(f)
	if err != nil {
		log.Panicln("Get absolute path error.")
	}
	p = strings.Replace(p, "\\", "/", -1)
	l := strings.LastIndex(p, "/") + 1
	return str.Substr(p, 0, l)
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsFile(f string) bool {
	b, err := os.Stat(f)
	if err != nil {
		return false
	}
	if b.IsDir() {
		return false
	}
	return true
}

func IsDir(p string) bool {
	b, err := os.Stat(p)
	if err != nil {
		return false
	}
	if b.IsDir() {
		return true
	}
	return false
}

//当前执行程序目录
func ExecutePath() string {
	dir, _ := os.Getwd()
	if len(os.Args) >= 2 {
		//服务启动时目录为参数传入进来
		dir = os.Args[1] + "\\"
	}
	return dir
}

//密码混淆
func EncryptPass(r string) (s string) {
	if r == "" {
		return ""
	} else {
		salt := "go2admin-haha" //加盐
		return encrypt.Md5(salt + r)
	}
}
