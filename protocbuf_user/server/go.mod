module github.com/example/user/server // 子模块定义‌:ml-citation{ref="7,8" data="citationList"}

go 1.20

require (
	github.com/example/user v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.60.1
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace github.com/example/user => ../ // 本地依赖替换‌:ml-citation{ref="7" data="citationList"}
