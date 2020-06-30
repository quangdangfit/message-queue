# Message Queue Service
### Golang, Echo, AMQP, MongDB

#### Setup
* Create config file: `cp config/config.sample.yaml config/config.yaml`
* Config database and amqp config
* Install require packages: `go mod vendor`

#### Startup
* Run Publisher: `go run -mod=vendor gomq/cmds/publisher/publisher.go`
* Run Consumer: `go run -mod=vendor gomq/cmds/consumer/consumer.go`
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
* `cmd/`: define commands
* `config/`: define configuration
* `dbs/`: init database connection, create index
* `packages/`: define packages
* `utils/`: common package

##### Package Structure
* `incoming/`: handle logic incoming messages (repo, model, handler)
* `inrouting/`: handle logic in routing key (repo, model)
* `outgoing/`: handle logic in routing key (repo, model, handler)
* `queue/`: contains publisher and consumer to send and receive messages
* `routers/`: define api url, request body, params
* `services/`: wrapper message before publish

#### 📙 Libraries
- [Echo Framework](https://echo.labstack.com/)
- [AMQP Package](https://godoc.org/github.com/streadway/amqp)

#### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
