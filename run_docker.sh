docker build  --no-cache -t iphone .
docker stop iphone_go
docker rm iphone_go
docker run -it --name iphone_go --net host iphone