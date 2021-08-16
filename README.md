# Teaching-school-load-balancing
该项目用于代理学校教务系统，通过检测教务系统各节点的健康度实现

## Todo
1. 服务监控(心跳)
2. 服务监听从localhost改全网段 √

## Service Port
- registry service -> 6000
- log service -> 6500
- redis service -> 7000
- tester service -> 7500
- discover service -> 8000

## scanner
教务系统扫描器，用于发现内网网段中隐藏的地址
项目结构
```shell
scanner
├── __init__.py # 为外部开发提供包
├── port.py # 端口扫描器（未使用）
├── scan.py # 扫描器主体
└── settings.py # 扫描器配置文件
```
