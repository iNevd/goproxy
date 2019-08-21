# 简介
构建goproxy， 收敛私有仓库访问 + 其他公共库访问


# 私有仓库方案

1. 生成`gitlab deploy key`， `xxx_id_rsa`
2. 在部署 `goproxy`机器上执行： 
   1. `git config --global url."git@Your Gitlab:".insteadOf "https://Your Gitlab/"`
   2. `export GIT_SSH_COMMAND='ssh -i FULL_PATH_DEPLOY_KEY -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no'`


# 使用

`./goproxy -cacheDir=cache -exclude=Your Gitlab -proxy=https://goproxy.io`

