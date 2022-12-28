APP=backend-api
LDFLAGS+="-s -w"

init:
	export GOPRIVATE=gitlab.com/*
	export GO111MODULE=on
	export GIT_TERMINAL_PROMPT=1
 	export GOPROXY=direct
	export GOSUMDB=off

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative utils/proto/*.proto

test: 
	go test -v -cover -covermode=atomic ./app_test/...

build:
	go build -v -ldflags $(LDFLAGS) -o $(APP) 


unittest:
	go test -short  ./internal/app_test

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

build-docker:
	docker build -t ${APP} .

run:
	docker-compose up --build -d

stop:
	docker-compose down

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint

build-staging:
	git checkout main
	git pull origin main
	git merge develop
	git push origin main
	git checkout develop