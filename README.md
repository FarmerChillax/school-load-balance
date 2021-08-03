# Teaching-school-load-balancing
该项目用于代理学校教务系统，通过检测教务系统各节点的健康度实现


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
