default: serve
APP?=grpclangserv
GOOS?=linux
GOARCH?=amd64
RELEASE?=0.0.1
CONTAINER_IMAGE?=colek42/${APP}
LS_GRPC_PORT?=4534


PROTO_GEN_GO := $(GOPATH)/bin/protoc-gen-go
DEP := $(GOPATH)/bin/dep
MEGACHECK := $(GOPATH)/bin/megacheck


PROTOC := $(shell which protoc)

ifeq ($(PROTOC),)
	PROTOC = must-rebuild
endif

UNAME := $(shell uname)

$(PROTOC):

ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

$(DEP):
	go get -u -v github.com/golang/dep/cmd/dep

$(MEGACHECK):
	go install -v honnef.co/go/tools/cmd/megacheck

compile: api/api.pb.go
api/api.pb.go: api/api.proto | $(PROTOC_GEN_GO) $(PROTOC)
	protoc --go_out=plugins=grpc:. api/api.proto


grpclangserv: $(shell find . -name '*.go')
	go build -o $(APP) github.com/colek42/$(APP)/cmd

dep:
	dep ensure -v -vendor-only

serve: dep compile lint test grpclangserv
	./$(APP)

test: $(shell find . -name '*.go')
	go test -cover -race ./...

lint: $(shell find . -name '*.go')
	$(MEGACHECK) ./...

clean:
	rm -rf ./vendor
	rm -rf $(APP)
	rm -rf ./api/api.pb.go
	rm -rf ./${APP}

build: clean compile dep lint test
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${APP}

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

run: container
	docker stop $(APP):$(RELEASE) | true && docker rm $(APP):$(RELEASE) || true
	docker run --name ${APP} -p ${LS_GRPC_PORT}:${LS_GRPC_PORT} --rm \
	-e "LS_GRPC_PORT=${LS_GRPC_PORT}" \
	-e "GOPATH"=/go \
	-v ${GOPATH}:/go \
	$(CONTAINER_IMAGE):$(RELEASE)



