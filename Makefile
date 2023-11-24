test:
	go test -race -coverprofile="coverage.out" -covermode=atomic ./...
	go tool cover -html="coverage.out"

lint:
	golangci-lint run

app_run:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d db
	sleep 15
	docker-compose -f docker-compose.yaml up -d init
	docker-compose -f docker-compose.yaml up -d migrate
	docker-compose -f docker-compose.yaml up -d uzum_delivery
#	docker-compose -f docker-compose.yaml up -d zookeeper
#	docker-compose -f docker-compose.yaml up -d kafka
#	docker exec -it uzum_delivery-kafka-1 kafka-topics --create --bootstrap-server localhost:29092 --topic my_topic --partitions 3 --replication-factor 1

app_run2:
	docker-compose -f docker-compose.yaml down -v
	docker-compose up -d
#Записать в Топик
#docker exec -it uzum_delivery-kafka-1 kafka-console-producer --broker-list localhost:29092 --topic my_topic

BIN:=$(CURDIR)/bin

install:
	GOBIN=$(BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2
	GOBIN=$(BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1


gen:
	protoc 	--proto_path=api/delivery_v1 \
			--proto_path=proto \
			--go_out=pkg/delivery_v1 --go_opt=paths=source_relative \
			--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
			--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
			--go-grpc_out=pkg/delivery_v1 --go-grpc_opt=paths=source_relative \
			--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway.exe \
			--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2.exe \
			--grpc-gateway_out=pkg/delivery_v1 --grpc-gateway_opt=paths=source_relative \
			--validate_out lang=go:pkg/delivery_v1 --validate_opt=paths=source_relative \
			--plugin=protoc-gen-validate=bin/protoc-gen-validate.exe \
			--openapiv2_out=allow_merge=true,merge_file_name=api_delivery_v1:docs \
			--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2.exe \
			api/delivery_v1/delivery.proto