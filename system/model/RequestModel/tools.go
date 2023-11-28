/*
 * @Date: 2022-04-12 09:11:39
 * @LastEditors: 春贰
 * @gitee: https://gitee.com/chun22222222
 * @github: https://github.com/chun222
 * @Desc:
 * @LastEditTime: 2022-11-10 16:58:58
 * @FilePath: \opcConnector\system\model\RequestModel\tools.go
 */
package RequestModel

type OpcServer struct {
	Name string `json:"name" binding:"required"`
}

type OpcTreePath struct {
	Path []string `json:"path" binding:"required"`
}

type OpcTags struct {
	Tags []string `json:"tags" binding:"required"`
}
