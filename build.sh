go build -o "./release/sd-webui-discord"
# 判断是否存在config.json
if [ ! -f "./release/config.json" ]; then
    echo "config.json not found, copy config.example.json to config.json"
    cp ./config.example.json ./release/config.json
fi