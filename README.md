# dig-the-grpc

## Benchtest

testing environment

* os        : ubuntu 16.04 LTS 64-bit
* processor : Intel® Core™ i7-8700 CPU @ 3.20GHz × 12 
* memory    : Kingston KHX2400C15/8G 2400 MHz DDR4 × 2
* graphics  : GeForce GTX 1060 6GB/PCIe/SSE2

> ONE-CONNECTION-PER-REQUEST test case

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    88.80ms  122.06ms   1.10s    89.20%
    Req/Sec     0.89k     0.89k    3.55k    88.00%
  35578 requests in 10.04s, 4.36MB read
Requests/sec:   3542.31
Transfer/sec:    444.34KB
```

> ONE-CONNECTION-PER-REQUEST test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   208.70ms    9.72ms 278.63ms   86.12%
    Req/Sec   239.35    131.19   480.00     59.02%
  9491 requests in 10.03s, 1.16MB read
Requests/sec:    946.72
Transfer/sec:    118.34KB
```

> ONLY-ONE-CONNECTION test case

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.77ms    1.20ms  17.15ms   73.57%
    Req/Sec    17.88k     1.71k   22.17k    76.00%
  712512 requests in 10.02s, 86.98MB read
Requests/sec:  71134.55
Transfer/sec:      8.68MB
```

> ONLY-ONE-CONNECTION test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   202.10ms    2.09ms 215.78ms   89.37%
    Req/Sec   245.17    102.05   484.00     72.08%
  9800 requests in 10.02s, 1.20MB read
Requests/sec:    978.32
Transfer/sec:    122.29KB
```

> CONNECTION-POOL-WITH-EXPANSION test case

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    47.05ms   97.37ms   1.03s    94.47%
    Req/Sec     2.12k     1.77k   10.52k    80.20%
  84426 requests in 10.08s, 10.33MB read
Requests/sec:   8372.39
Transfer/sec:      1.02MB
```

> CONNECTION-POOL-WITH-EXPANSION test case (server sleep for 200ms before respose)

```bash
$ wrk -t4 -c200 -d10s http://localhost:18080/performance

Running 10s test @ http://localhost:18080/performance
  4 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   202.45ms    7.56ms 273.20ms   97.47%
    Req/Sec   246.81     95.41   494.00     67.16%
  9800 requests in 10.03s, 1.20MB read
Requests/sec:    977.39
Transfer/sec:    122.17KB
```

## Reference

* [grpc](https://grpc.io/)
* [grpc-go-pool](https://github.com/processout/grpc-go-pool)
* [pooling grpc connections](https://mycodesmells.com/post/pooling-grpc-connections)
