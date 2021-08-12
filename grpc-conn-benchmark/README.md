## environment

```text
OS     : Ubuntu 18.04.5 LTS x86-64
Kernel : Linux 4.15.0-106-generic
CPU    : Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz × 12 
Memory : Kingston KHX2400C15/8G 2400 MHz DDR4 × 2
GPU    : GeForce GTX 1060 6GB/PCIe/SSE2
```

### "connect: cannot assign requested address"解决方法

netstat | grep TIME_WAIT

sudo vim /etc/sysctl.conf

```text
# ls -l /proc/sys/net/ipv4/
# 修改tcp/ip协议配置, 修改默认的TIME_WAIT TIMEOUT时间, 默认为60s
sysctl -w net.ipv4.tcp_fin_timeout=30
# 修改tcp/ip协议配置, 释放TIME_WAIT端口给新连接使用, 默认为0
sysctl -w net.ipv4.tcp_tw_reuse = 1
# 修改tcp/ip协议配置, 快速回收socket资源, 默认为0
sysctl -w net.ipv4.tcp_tw_recycle=1
```

sudo sysctl -p

### ONE_CONNECTION_PER_REQUEST test case

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```

### ONE_CONNECTION_PER_REQUEST test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```

### ONLY_ONE_CONNECTION test case

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```

### ONLY_ONE_CONNECTION test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```

### CONNECTION_POOL_WITH_EXPANSION test case

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```

### CONNECTION_POOL_WITH_EXPANSION test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c100 -d10s http://localhost:18888/performance
```
