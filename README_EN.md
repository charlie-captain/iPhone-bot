# iPhone Bot

iPhone Pro Telegram Bot, Default city is GuangZhou.

## Usage

1. settings.json edit

```
{
  # telegram bot TOKEN
  "bot_token": "",
  
  # chat id
  "chat_id": 00000000,
  
  # store id, search from store.json
  "store_list": [
    "R577",
    "R639"
  ],
  
  # model id, from website http request
  "model_list": [
    
  ],
  
  # fetch duration: 3s/1min/1h
  "fetch_duration": "3s",
  
  # proxy:port(http://127.0.0.1:7890), default system proxy
  "proxy": "",
  
  # region(default china), example:: "https://www.apple.com/hk/", so input "hk"
  "region": "",
}
```

2. Run Program

```
windows double click iphoneBot.exe

linux/Mac run iphoneBot in terminal, maybe chmod +x ./iphoneBot
```

3. Docker

```
chmod +777 ./run_docker.sh 
./run_docker.sh
```

## Dev

1. clone project
2. fix bug
3. run command/build.sh

