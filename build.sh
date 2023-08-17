###
 # @Author: SpenserCai
 # @Date: 2023-08-17 11:04:55
 # @version: 
 # @LastEditors: SpenserCai
 # @LastEditTime: 2023-08-17 23:28:05
 # @Description: file content
### 
export GOOS=linux
go build -o "./release/sd-webui-discord"

export GOOS=windows
go build -o "./release/sd-webui-discord.exe"
# 判断是否存在config.json
if [ ! -f "./release/config.json" ]; then
    echo "config.json not found, copy config.example.json to config.json"
    cp ./config.example.json ./release/config.json
fi