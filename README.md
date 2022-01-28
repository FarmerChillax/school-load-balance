# School-load-balancing
该项目用于代理学校教务系统，通过检测教务系统各节点的健康度实现

## Todo
1. ~~服务监控(心跳)~~ √
2. ~~服务监听从localhost改全网段~~ √

## Service Port
- registry service -> 6000
- log service -> 6500
- redis service -> 7000
- tester service -> 7500
- discover service -> 8000

## How to make your service
1. 向服务中心注册
2. 接受注册中心服务的更新
3. 编写你自己服务的逻辑

## cmd
存储各个服务的启动项，每个服务的`port`存放端口 `host`存放部署主机的外部联系地址(`domain name` or `ip`)

## Services
各个服务的功能介绍

### discover
教务系统扫描器，用于发现内网网段中隐藏的地址

### tester
测试redis中的数据

### log
日志服务，存储系统日志

### storage
存储服务，目前使用`redis`

