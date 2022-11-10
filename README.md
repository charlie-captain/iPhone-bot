# iPhone Bot

[English](./README_EN.md)

iPhone Pro Telegram Bot, 默认城市 GuangZhou.

## Usage

1. settings.json 编辑

```
{
  # 机器人TOKEN
  "bot_token": "",
  
  # 聊天id, 自己的id或者群聊
  "chat_id": 00000000,
  
  # 店铺代号, 从根目录store.json查找, 最好是临近的店铺
  "store_list": [
    "R577",
    "R639"
  ],
  
  # 监听的型号代号, 从iPhone购买页面地址获取, 默认空则为 Pro + Pro Max(不建议超过3个)
  "model_list": [
    
  ],
  
  # 查询间隔时间 3s/1m/1h
  "fetch_duration": "3s",
  
  # 代理地址:端口(例 http://127.0.0.1:7890), 默认为系统代理
  "proxy": "",
  
  # 地区(默认中国内地), 例如 中国香港 https://www.apple.com/hk/ 输入 hk 即可, 后续会改成机器人配置
  "region": "",
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

## Dev

1. clone 项目
2. 改动代码
3. 运行 command/build.sh

## 声明

本项目仅供学习使用, 不能用于商业非法等用途, 本人不承担任何法律责任