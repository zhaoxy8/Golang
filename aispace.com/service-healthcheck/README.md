### 一、Network True 随机端口服务可用性检测
```
容器读取环境变量配置
export PROJECT_ENV=M2-STG
export PROJECT_THRESHOLD=1
export NACOS_CLUSTER=https://m-nacos-stg.bmw-emall.cn
export DINGDING_SECRET=SECf256c8ae405df3f958b613aa41e723b3a3aeaddff182e69547288dda2fa31a66
export DINGDING_URLADDRESS=https://oapi.dingtalk.com/robot/send?access_token=b291d1fc6728b6be13b786ec44809bb36425a33cec022c7782c98d19356f6f87

Ecom 监控报警新增规则，每2分钟检查1次：
1.如果服务可用节点数小于等于1报警
2.检查服务节点健康状态"/test/health"访问不可达时报警
3.服务端口连接拒绝的时候报警


M2监控报警新增规则，每2分钟检查1次：
1.如果服务可用节点数小于等于1报警
2.检查服务节点健康状态"/publicApi/health"访问不通时报警
3.服务端口连接拒绝的时候报警
```
