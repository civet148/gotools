1.github上下载一个cpp包：https://github.com/google/protobuf/releases make make install安装即可

2.protoc-gen-go
go get -u github.com/golang/protobuf/protoc-gen-go

3.安装protoc-gen-micro
go get github.com/micro/protoc-gen-micro

4. 启动consul
* Windows

```batch
@echo off
consul agent -server -bootstrap-expect 1 -data-dir .\data -node=s1 -bind=127.0.0.1 -http-port=8500  -dns-port=5600 -rejoin -config-dir=.\etc\consul.d -client 0.0.0.0
pause
```

* Linux

```shell script
nohup consul agent -server -bootstrap-expect 1 -data-dir /data/consul/ -node=s1 -bind=127.0.0.1 -http-port=8500  -dns-port=8600 -rejoin -config-dir=/etc/consul.d/ -client 0.0.0.0 > /var/log/consul.log 2>&1 &
```