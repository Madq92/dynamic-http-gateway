#!/bin/bash

# 定义环境变量
export LANG="en_US.UTF-8"

# 业务配置路径
CC_BASE_URL=${CC_BASE_URL:-cc.56qq.cn}
gitpath=$GIT_PATH

printMsg(){
  echo  "$(date +'%Y-%m-%d %H:%M:%S')-$1"
}

if [ "${gitpath}" == "" ]; then
    echo "无CC参数，使用本地配置项"
else
    CC_URL="http://${CC_BASE_URL}/business/pullZipFile.do?pathId=${gitpath}&type=3"
    if [ "${CC_USERNAME}" != "" ]; then
        CC_URL="${CC_URL}&userName=${CC_USERNAME}"
    fi
    if [ "${CC_PASSWORD}" != "" ]; then
        CC_URL="${CC_URL}&password=${CC_PASSWORD}"
    fi
    if [ "${DMS_API_KEY}" != "" ]; then
        CC_URL="${CC_URL}&apiKey=${DMS_API_KEY}"
    fi

    echo "CC_URL=========: ${CC_URL}"

    TARBALL=/tmp/cc.tar.gz
    wget -v -t 5 ${CC_URL} -O ${TARBALL}


    tar zxfv ${TARBALL} -C /app --strip-components=1

    if [ $? != 0 ]; then
        printMsg '----------------拉取配置文件失败------------'
        exit 1
    fi
fi


# 启动服务
mkdir -p /var/logs/ulog
nohup ./dynamic-http-gateway > /var/logs/ulog/startup.log 2>&1 & tail -f /var/logs/ulog/startup.log
