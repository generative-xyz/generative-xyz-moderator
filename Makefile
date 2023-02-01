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

build-gpu-develop:
	git checkout gpu-develop && git merge develop && git push origin gpu-develop && git checkout develop

build-staging:
	git checkout staging && git pull origin staging && git merge develop && git push origin staging
	git checkout gpu-staging && git pull origin gpu-staging && git merge staging && git push origin gpu-staging
	git checkout develop

start-docker:
	docker-compose stop api_service && docker-compose up -d api_service && docker logs -f api_service
	
reload-docker:
	docker-compose stop api_service && docker-compose build api_service  && docker-compose up -d api_service && docker logs -f api_service

exec-docker:
	docker exec -it api_service bash 

start-xserver:
	docker-compose stop api_xserver && docker-compose build api_xserver && docker-compose up -d api_xserver