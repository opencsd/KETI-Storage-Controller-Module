export PATH="$PATH:$(go env GOPATH)/bin" #protoc 컴파일러가 플러그인 찾을 수 있게 설정 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative config/config.proto 