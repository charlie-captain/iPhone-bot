rm -f ./logs.txt
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
echo "build success"
git fetch origin
git add .
time=$(date "+%Y-%m-%d %H:%M:%S")
echo $time
git commit -m "bump : ${time}"
git rebase origin/main
git push origin main