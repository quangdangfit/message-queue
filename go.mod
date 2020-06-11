module gomq

go 1.14

require (
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/manucorporat/try v0.0.0-20170609134256-2a0c6b941d52
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
	github.com/valyala/fasttemplate v1.1.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	transport/lib v0.0.0-20200609030708-5cbccf123a48
)

replace transport/lib => gitlab.com/transport4/lib.git v0.0.0-20200609030708-5cbccf123a48
