@echo off

echo generating...

protoc -I=.  --proto_path=%GOPATH%\src  --micro_out=.  --go_out=.   *.proto

echo generate over

pause
