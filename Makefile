RANDOM_OBJECT_OUT := "randonObjectID.bin"
CONTROL_OUT := "control.bin"
API_OUT := "api/api.pb.go"
PKG := "github.com/AntanasMaziliauskas/grpc/cmd"
RANDOM_OBJECT_PKG_BUILD := "${PKG}/utils/randomObjectID"
CONTROL_PKG_BUILD := "${PKG}/control"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all api build_server_docker build_node_docker build_randomObjectID build_control install_control

all: build_server_docker build_node_docker build_randomObjectID install_randomObjectID build_control install_control
 
api/api.pb.go: api/api.proto
	@protoc -I api/ \
		-I${GOPATH}/src \
		--go_out=plugins=grpc:api \
		api/api.proto

api: api/api.pb.go ## Auto-generate grpc go sources

dep: ## Get the dependencies
	@go get -v -d ./...

build_randomObjectID: 
	@go build -i -v -o $(RANDOM_OBJECT_OUT) $(RANDOM_OBJECT_PKG_BUILD)

install_randomObjectID:
	@go install -i $(RANDOM_OBJECT_PKG_BUILD)

build_control: dep api ## Build the binary file for server
	@go build -i -v -o $(CONTROL_OUT) $(CONTROL_PKG_BUILD)

install_control:
	@go install -i $(CONTROL_PKG_BUILD)

build_server_docker: dep api ## Build the binary file for server
	@docker build -t server -f Dockerfile.server .

build_node_docker: dep api ## Build the binary file for node
	@docker build -t node .

clean: ## Remove previous builds
	@rm $(SERVER_OUT) $(CLIENT_OUT) $(API_OUT)

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
