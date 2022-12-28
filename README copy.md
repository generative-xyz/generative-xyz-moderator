# OG device APIs
## Setting up

- Please create .env file with the following content.

    ```
        "swagger_domain" : "{your_server_host}:{your_server_port}",
        "PORT": {your_server_port},
        "DEBUG": {bool: true | false},
        "CONTEXT_TIMEOUT": 30,
        "JAEGER_SERVICE_NAME" : "og-base-project",
        "JAEGER_REPORTER_LOG_SPANS": "true",
        "JAEGER_AGENT_HOST": "127.0.0.1",
        "JAEGER_AGENT_PORT": "6831",
        "AUTH_SECRET_KEY": "{create_a_service_key}",
        "MONGO_USER": "{your_mongo_db_user}",
        "MONGO_PASSWORD": "{your_mongo_db_password}",
        "MONGO_HOST": "{your_mongo_db_host}",
        "MONGO_PORT": "{your_mongo_db_port}",
        "MONGO_DB": "{your_mongo_db_name}",
    ```

- After config file has been created, please use the command below to get all neccessary modules.

    ```
        go mod tidy
    ```

- Done!
- We can start our application via Vscode or anything to develop.

## Development:
- Please create a launch.json via VSCode with the content below.

    ```
        {
            // Use IntelliSense to learn about possible attributes.
            // Hover to view descriptions of existing attributes.
            // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
                "version": "0.2.0",
                "configurations": [
                    {
                        "name": "Launch Package",
                        "type": "go",
                        "request": "launch",
                        "mode": "debug",
                        "program": "${workspaceFolder}",
                        "env": {
                            "swagger_domain" : "{your_server_host}:{your_server_port}",
                            "PORT": {your_server_port},
                            "DEBUG": {bool: true | false},
                            "CONTEXT_TIMEOUT": 30,
                            "JAEGER_SERVICE_NAME" : "og-base-project",
                            "JAEGER_REPORTER_LOG_SPANS": "true",
                            "JAEGER_AGENT_HOST": "127.0.0.1",
                            "JAEGER_AGENT_PORT": "6831",
                            "AUTH_SECRET_KEY": "{create_a_service_key}",
                            "MONGO_USER": "{your_mongo_db_user}",
                            "MONGO_PASSWORD": "{your_mongo_db_password}",
                            "MONGO_HOST": "{your_mongo_db_host}",
                            "MONGO_PORT": "{your_mongo_db_port}",
                            "MONGO_DB": "{your_mongo_db_name}",
                        
                        },
                        "args": ["server"]
                    }
                ]
            }
    ```

- Note: we will use **env** field in the **launch.json** instead of **.env**, if VSCode is used.

## Commands

- All commands are defined in **Makefile** like the following list:
    - **make test**, is used to run unit test.
    - **make build-docker**, is used to build docker image.
    - **make run**, is used to run the built image.
    - **make lint-prepare**, is used to install lint.
    - **make lint**, is used to run lint test.

## APIs 
- Please try out APIs with this [APIs documentation](https://dev-hybrid.autonomous.ai/credit-api-docs/index.html#/), and don't forget to fill out the **Api-Key** by the **Authorizations** button :D.

- Application response codes and messages are defined in this application:

    ```go
        const (
            //Success has prefix with 8"0"
            Success      = 1
            
            //Error has prefix with 8"1"
            Error = -1
        )

        //Message
        var ResponseMessage = map[int]struct {
            Code    int
            Message string
        }{
            Success:{Success, "Success"},
	        Error: {Error, "Failed."},
        }

    ```

## Tech stack using
- [Gorm](https://gorm.io/docs/)
- [Mux](https://github.com/gorilla/mux)
- [Zap logger](https://github.com/uber-go/zap)
- [Clean Architecture](https://github.com/bxcodec/go-clean-arch)
- [Swagger](https://github.com/swaggo/swag#general-api-info)
- [Swagger with Mux in Golang](https://devopstutorial.tech/swagger-ui-setup-for-go-rest-api-using-swaggo/)
