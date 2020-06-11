module gomq

go 1.14

require (
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/viper v1.7.0
	github.com/streadway/amqp v1.0.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	transport/lib v0.0.0-20200609030708-5cbccf123a48
)

replace transport/lib => gitlab.com/transport4/lib.git v0.0.0-20200609030708-5cbccf123a48