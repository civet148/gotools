1.github上下载一个cpp包：https://github.com/google/protobuf/releases make make install安装即可

```shell script
$ wget https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protobuf-cpp-3.13.0.tar.gz
$ tar xvfz protobuf-cpp-3.13.0.tar.gz
$ cd protobuf-cpp-3.13.0
$ ./configure && make && make install
```
2. 安装go插件 protoc-gen-go

```shell script
$ go get -u github.com/golang/protobuf/protoc-gen-go
```


3.安装go-micro插件 protoc-gen-micro
```shell script
$ go get github.com/micro/protoc-gen-micro
```

4. 启动consul
* Windows

```batch
@echo off
consul agent -server -bootstrap-expect 1 -data-dir .\data -node=s1 -bind=127.0.0.1 -http-port=8500  -dns-port=5600 -rejoin -config-dir=.\etc\consul.d -client 0.0.0.0
pause
```

* Linux

```shell script
$ sudo mkdir -p /data/consul
$ sudo chmod 0777 /data/consul
nohup consul agent -server -bootstrap-expect 1 -data-dir /data/consul -node=s1 -bind=127.0.0.1 -http-port=8500  -dns-port=8600 -rejoin -config-dir=/etc/consul.d/ -client 0.0.0.0 > /tmp/consul.log 2>&1 &
```