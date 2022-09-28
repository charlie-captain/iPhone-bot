# iPhone Bot

iPhone Pro Telegram Bot, 默认城市 GuangZhou.

## Usage

1. settings.json 编辑

```
{
  # 机器人TOKEN
  "bot_token": "",
  
  # 聊天id
  "chat_id": 00000000,
  
  # 店铺代号, 从根目录store.json查找, 最好是临近的店铺
  "store_list": [
    "R577",
    "R639"
  ],
  
  # 查询间隔时间 3s/1min/1h
  "fetch_duration": "3s",
  
  # 代理地址:端口(例 http://127.0.0.1:7890), 默认为系统代理
  "proxy": ""
}
```

2. 本地运行程序

```
windows 运行 iphoneBot.exe

linux/Mac 运行 iphoneBot, 需要命令行设置权限 chmod +x ./iphoneBot
```

3. Docker部署

```
//命令行运行
chmod +777 ./run_docker.sh 
./run_docker.sh
```
