### protobuf generation with protoc-gen-micro and protoc-gen-go
protoc -I=.  --proto_path=$GOPATH/src  --micro_out=.  --go_out=.   *.proto
