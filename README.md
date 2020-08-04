# Message Queue Service
### Golang, Echo, AMQP, MongDB

#### Setup
* Create config file: `cp config/config.sample.yaml config/config.yaml`
* Config database and amqp config
config/config.yaml
```yaml
...
mode: # 0: run publisher and consumer, 1: run publisher, 2: run consumer 
...
```

* Install require packages: `go mod vendor`

#### Startup
* Run: `go run -mod=vendor gomq/main.go`
* Document at: http://localhost:8080/swagger/index.html

![](https://i.imgur.com/4qewM7a.png)

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

| Fields       | Type          | Required | Not Null | Description                       |
|:-------------|:-------------:|:--------:|:--------:|:----------------------------------|
| routing_key  | string        | YES      | YES      | Routing key                       |
| payload      | json          | YES      | YES      | Message content (json)            |
| origin_model | string        | NO       | NO       | Object model                      |
| origin_code  | string        | NO       | NO       | Object code                       |


#### Diagram
![alt text](https://imgur.com/NXuvQLG.jpg "Repository Pattern")


#### Structure
```
├── app  
│   ├── api             # Handle request & response
│   ├── dbs             # Database Layer
│   ├── models          # Models
│   ├── queue           # AMQP Layer
│   ├── repositories    # Repositories Layer
│   ├── router  
│   │   └── v1          # Router api v1  
│   ├── schema          # Sechemas  
│   ├── services        # Business Logic Layer  
│   └── utils           # Utilities  
├── config              # Config's files  
```

#### 📙 Libraries
- [Echo Framework](https://echo.labstack.com/)
- [AMQP Package](https://godoc.org/github.com/streadway/amqp)

#### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
