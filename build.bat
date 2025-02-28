@echo off
chcp 65001
echo ===== 开始构建 OPC Connector =====

echo.
echo ===== 设置 GOPROXY =====
set GOPROXY=https://goproxy.cn,direct

echo.
echo ===== 检查 Go 环境 =====
go version

echo.
echo ===== 检查依赖 =====
go list -m all

echo.
echo ===== 更新依赖 =====
go mod tidy
if %ERRORLEVEL% NEQ 0 (
    echo 依赖更新失败
    pause
    exit /b 1
)

echo.
echo ===== 检查代码 =====
go vet ./...
if %ERRORLEVEL% NEQ 0 (
    echo 代码检查失败
    pause
    exit /b 1
)

echo.
echo ===== 构建 32 位版本 =====
set GOARCH=386
go build -v -ldflags="-s -w" -o opcConnector-win32.exe
if %ERRORLEVEL% EQU 0 (
    echo 32位版本构建成功：opcConnector-win32.exe
) else (
    echo 32位版本构建失败
    pause
    exit /b 1
)

echo.
echo ===== 构建 64 位版本 =====
set GOARCH=amd64
go build -v -ldflags="-s -w" -o opcConnector-win64.exe
if %ERRORLEVEL% EQU 0 (
    echo 64位版本构建成功：opcConnector-win64.exe
) else (
    echo 64位版本构建失败
    pause
    exit /b 1
)

echo.
echo ===== 构建完成 =====
dir *.exe

pause 