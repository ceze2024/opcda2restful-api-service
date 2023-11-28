/*
 * @Date: 2022-02-14 14:09:57
 * @LastEditors: 春贰
 * @Desc:
 * @LastEditTime: 2023-07-11 13:53:36
 * @FilePath: \opcConnector\system\common\initial\config.go
 */
package initial

const ConfigToml = `[config]
[config.App]
Secret = 'ADSDWW1DSADSADSAWJJK'
HttpPort = 9022
OpcHost = '127.0.0.1'
OpcServer = 'Kepware.KEPServerEX.V6'
KeepConn = 1

[config.Zaplog]
Director = 'runtime/log'
Level = 'debug' 
`
