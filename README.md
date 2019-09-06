# 支撑http业务的动态路由网关

### 8080端口

业务端口，负责代理业务的请求，根据不同的规则代理到对应的Target（具体的业务Service）

### 8081端口

管理端口，负责动态更新规则，更新规则接口：

~~~
POST /rule  更新规则，全量

body:
[
  {
    "name": "1",
    "targets": [
      "http://localhost:9001"
    ],
    "matchRules": [
      {
        "name": "headersEq",
        "param": [
          "key1",
          "value1"
        ]
      }
    ]
  },
  {
    "name": "2",
    "targets": [
      "http://localhost:9000"
    ]
  }
]


GET /rule  获取当前规则
~~~

##### body参数说明
name：该条规则的名字，修改时当ID用

targets：目标服务，可以多个，如果规则符合，那么多个targets之间是roundRobin

matchRules：匹配的规则，没有规则就是全匹配，多条规则则每条规则都要符合

#### matchRules的定义

Matcher | Description | Example
---|---|---
HeadersEq | header等于匹配 | name=HeadersEq,param=["key","value","key2","value2"]
HeadersRegex | header正则匹配 | name=HeadersRegex,param=["key","valueRegex","key2","value2Regex"]
HeadersHas | header包含 | name=HeadersHas,param=["key","key2"]
Path | path全匹配 | name=Path,param=["/path1","/path2"]
PathPrefix | path前缀匹配 | name=PathPrefix,param=["/pathPrefix"]
Method | 方法匹配 | name=Method,param=["GET","POST"]
Query | query参数 | name=Query,param=["key","value","key2","value2"]
Body | body参数 | name=Body,param=["key","value","key2","value2"]

#### TODO list

路由规则信息存ETCD

支持水平扩容