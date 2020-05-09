# 日志监控收集系统 <br>
- logAgent收集日志到kakfa中，logTransfer负责把日志存储到ES存储中 <br>
- logAgent通过ETCD读取配置文件信息，更新配置. <br>
- 主要使用第三方库 <br>
```
"gopkg.in/ini.v1"
"github.com/Shopify/sarama"
"github.com/hpcloud/tail"
```

