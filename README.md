# Geo
> 收集nginx日志信息，对日志信息**IP地址**过滤分析

## Flow
```text
                                                                +------------------------+
                        +--------+                              |                        |
+-----------+           |        |          +----------+        |                        |
|fluent-bit |---------->| kafka  |<---------|  golang  |------->|     elasticsearch      |
+-----------+           |        |          +----------+        |                        |
                        +--------+                              |                        |
                                                                +------------------------+
```

## 入门开始

### Nginx日志输出配置
```nginx
http {
        # Logging Settings
        log_format main   '{"@timestamp":"$time_iso8601",'
                        '"@source":"$server_addr",'
                        '"hostname":"$hostname",'
                        '"ip":"$http_x_forwarded_for",'
                        '"client":"$remote_addr",'
                        '"request_method":"$request_method",'
                        '"scheme":"$scheme",'
                        '"domain":"$server_name",'
                        '"referer":"$http_referer",'
                        '"request":"$request_uri",'
                        '"args":"$args",'
                        '"size":$body_bytes_sent,'
                        '"status": $status,'
                        '"responsetime":$request_time,'
                        '"upstreamtime":"$upstream_response_time",'
                        '"upstreamaddr":"$upstream_addr",'
                        '"http_user_agent":"$http_user_agent",'
                        '"https":"$https"'
                        '}';

        access_log /var/log/nginx/access.log main;
        error_log /var/log/nginx/error.log;
}
```

### Fluent-bit配置
```ini
[INPUT]
    name              tail
    path              /var/log/nginx/*.log
    DB                /fluent-bit/etc/fluent-bit.db
    Mem_Buf_Limit     5MB
    Skip_Long_Lines   On
    Refresh_Interval  10
    Parser            json

[OUTPUT]
    Name        kafka
    Match       *
    Brokers     192.168.52.172:9092
    Topics      nginx-log
```

### 安装和运行

```bash
./geo
```

#### 配置说明

```yaml
kafka:
  host: "kafka:9092"                        # kafka服务器ip，多个kafka以逗号隔开 例如："kafka1:9092,kafka2:9092,kafka3:9092"
  topic: "nginx-log"                        # 消费topic
  groupId: "1"                              # kafka groupid
  numPartitions: 3                          # 自动创建kafka分区数量
  replication: -1                           # 自动创建kafka副本数量
es:
  host: http://es:9200                      # es地址
  preName: logstash-                        # es文档前缀

chunzhen:                          
  path: config/qqwry.dat                    # 纯真离线库路径
ip2location:
  path: config/IP2LOCATION-LITE-DB11.BIN    # ip2location离线数据库路径
```

## 源码安装
```bash
go mod download
go build
```