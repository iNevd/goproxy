# 简介

构建goproxy， 收敛私有仓库访问 + 其他公共库访问

```
                          +----------------------------------> go get(local)
                          |
                     is private
                          |
                      +---+---+           +---------+
go get xxx  +-------> |goproxy| +-------> |yyy porxy| +---> http(remote)
                      +-------+           +---------+
```

# 私有仓库方案

1. 配置好可以直接`clone`私有库代码的环境（id_rsa加到gitlab/github）
2. 在`goproxy`程序所在目录放置 `.gitconfig`文件（注意替换 `git.xx.com`）
   1. 将 https形式访问私有库自动替换为使用ssh形式

      ```bash
      [url "git@git.xxxx.com:"]
      	insteadof = https://git.xxx.com/
      ```
3. 执行脚本，设置好环境变量（注意替换 `git.xxx.com`)

```bash
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
if [ "x$SHELL" != 'x/bin/bash' ]
then
    BASEDIR="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
fi

export HOME=${BASEDIR}
export GOPROXY=goproxy.cn
export GOPRIVATE=git.xxx.com
${BASEDIR}/goproxy -addr=0.0.0.0:8081
```

# 其他使用

如需将代理数据缓存，避免每一次请求远程， 可以使用 `-cacheDir`指定缓存目录
