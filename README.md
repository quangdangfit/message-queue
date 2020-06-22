# Message Queue Service
### Golang, Echo, AMQP

#### Setup
* Create config file: `cp config-sample.yaml config.yaml`
* Config database and amqp config
* Install require packages: `go mod vendor`

#### Startup
* Run Publisher: `go run -mod=vendor gomq/cmd/publisher/mqPublisher.go`
* Run Consumer: `go run -mod=vendor gomq/cmd/consumer/mqConsumer.go`
* Publish message:
```
curl --location --request POST 'localhost:8080/api/v1/queue/messages' \
--header 'Content-Type: application/json' \
--data-raw '{
    "routing_key": "routing.key",
    "payload": {
        "name": "This is message"
    },
    "origin_code": "CODE",
    "origin_model": "MODEL"
}'
```

#### Diagram
![alt text](https://imgur.com/NXuvQLG.jpg "Repository Pattern")


#### Structure
* `routers/`: define api url, request body, params
* `services/`: wrapper message before publish
* `msgQueue/`: contains publisher and consumer to send and receive messages
* `msgHandler/`: handle messages (store DB, call api,...)
* `models/`: define orm models
* `config/`: define configuration

#### 📙 Libraries
- Echo: https://echo.labstack.com/
- AMQP Package: https://godoc.org/github.com/streadway/amqp

#### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
