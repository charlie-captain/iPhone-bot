# iPhone-bot
iPhone 14 Pro Telegram Bot, 默认城市 GuangZhou.

## Usage

1. .env edit
```
# 机器人TOKEN
BOT_TOKEN=xxx

# 自己的id
MY_ID=xxx

# Url 监听地址, 需要修改store后面的代号, product的机型代号, 可以在网页端抓包获取
URL=https://www.apple.com.cn/shop/pickup-message-recommendations?mt=compact&searchNearby=true&store=R639&product=MPWL3CH/A

# 查询间隔时间 3s/1min/1h
FETCH_DURATION=3s
```

2. run program


## Q & A

### store 是商店的代号，product 是机型，如何获取？

1. 打开购买链接 https://www.apple.com.cn/shop/buy-iphone/iphone-14-pro
2. 选择好对应的颜色，容量，到最后添加购物车的步骤
3. 按F12打开开发者控制台，选择 network(网络) -> Fetch/XHR
4. 网页鼠标点击查看其他零售店，选择好地区后，查看控制台会多一条 pickup-message-recommendations
5. 直接将 RequestURL 复制到上面的环境中，运行即可

