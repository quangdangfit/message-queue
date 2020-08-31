# Message Queue Service
### Golang, Gin, AMQP, MongDB

### Setup
* Create config file: `cp config/config.sample.yaml config/config.yaml`
* Config database and amqp config
config/config.yaml
```yaml
...
mode: # 0: run publisher and consumer, 1: run publisher, 2: run consumer 
...
```

* Install require packages: `go mod vendor`

### Startup
* Run: `go run -mod=vendor main.go`
* Document at: http://localhost:8080/swagger/index.html

![](https://i.imgur.com/Eh1KZAK.png)

### Publish message:
* **REST**:
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

* **RPC**:

Service support rpc for publishing, create the client as below to call rpc:
```go
package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	body := map[string]interface{}{
		"routing_key": "routing.key",
		"payload": map[string]interface{}{
			"name": "This is message",
		},
		"origin_code":  "CODE",
		"origin_model": "MODEL",
	}
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Call("OutRPC.Publish", body, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}
```

* **Body**:

| Fields       | Type          | Required | Not Null | Description                       |
|:-------------|:-------------:|:--------:|:--------:|:----------------------------------|
| routing_key  | string        | YES      | YES      | Routing key                       |
| payload      | json          | YES      | YES      | Message content (json)            |
| origin_model | string        | NO       | NO       | Object model                      |
| origin_code  | string        | NO       | NO       | Object code                       |

### Diagram
![alt text](https://i.imgur.com/KwUNR1V.png)


### Structure
```
├── app  
│   ├── api             # Handle request & response
│   ├── dbs             # Database Layer
│   ├── models          # Models
│   ├── queue           # AMQP Layer
│   ├── repositories    # Repositories Layer
│   ├── router          # Router api v1  
│   ├── schema          # Sechemas  
│   ├── services        # Business Logic Layer  
├── config              # Config's files 
├── docs                # Swagger API document
├── pkg                 # Common packages
│   ├── app             # App's packages
│   └── utils           # Utilities
```

### 📙 Libraries
- [Gin](https://godoc.org/github.com/gin-gonic/gin)
- [AMQP](https://godoc.org/github.com/streadway/amqp)

### Contributing
If you want to contribute to this boilerplate, clone the repository and just start making pull requests.
