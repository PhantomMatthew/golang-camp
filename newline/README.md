# falcon
11
1. 进行数据库建立，在根目录执行

   `go run src/cmd/migrate.go --config path:./src/srv/config/dev.json`

2. 启动customer服务，在根目录执行

   `go run src/srv/customer/main.go --config ../config/staging.json`

3. 启动api服务，在根目录执行

   `go run src/api/customer_aggregation/main.go --config ../config/staging.json`

4. 启动网关，在根目录执行

   `go run src/gateway/central/main.go api --handler=http`



进行测试，可以使用postman执行http://localhost:8080/customer/users。



建立一个新的service，使用`micro new {name}`进行建立，然后删除此service的go.mod，使用根目录的go.mod。   
另外在handler/customer.go subscriber/customer.go 还有main.go中修改导入包的路径，由于整个项目使用的是newline.com/newline，需要对应相对的目录结构进行上述三个文件的修改。



以下为建立新service的memo记录，建立后需要使用protoc命令生成相应protobuf的处理Go代码。

Creating service go.micro.service.customer in customer  

.   
├── main.go   
├── generate.go   
├── plugin.go   
├── handler  
│   └── customer.go  
├── subscriber  
│   └── customer.go  
├── proto/customer   
│   └── customer.proto   
├── Dockerfile   
├── Makefile   
├── README.md   
├── .gitignore    
└── go.mod   


download protobuf for micro:   

`brew install protobuf`
`go get -u github.com/golang/protobuf/{proto,protoc-gen-go}`
`go get -u github.com/micro/protoc-gen-micro/v2`

compile the proto file customer.proto:

`cd customer`   
`protoc --proto_path=.:$GOPATH/src --go_out=. --micro_out=. proto/customer/customer.proto`
去除omitempty
`ls proto/customer/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'`
