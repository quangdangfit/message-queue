package queue

type AMQPConfig struct {
	AMQPUrl      string
	Host         string
	Port         string
	Vhost        string
	Username     string
	Password     string
	ExchangeName string
	ExchangeType string
	QueueName    string
}
