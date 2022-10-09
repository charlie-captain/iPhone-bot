pwd
git fetch origin
git add .
time=$(date "+%Y-%m-%d %H:%M:%S")
echo $time
git commit -m "bump : ${time}"
git rebase origin/main
git push origin main