package rabbitmq

import (
	"fmt"
	"post_service/config"
	"post_service/pkg/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
) 

var(
	l logger.Logger
	cfg config.Config
)

type Config struct {
	URL string
	WaitTime time.Duration
	Attempts int
}

type Connection struct {
	ConsumerExchange string
	Config
	Connection *amqp.Connection
	Channel *amqp.Channel
	Delivery <-chan amqp.Delivery
	
}

func init(){
	cfg = *config.NewConfig()
	l = *logger.New(cfg.App.LogLevel)
}

func New(consumerExchange string, cfg Config )  *Connection {
	conn:=&Connection{
		ConsumerExchange: consumerExchange,
		Config: cfg,
	}
	return conn
}

func failOnError(err error, msg string){
	if(err!=nil){
		l.Fatal(err,msg)
	}
}

func (c *Connection)AttemptConnect() error{
	var err error
	for i:=c.Attempts;i>0;i--{
		if err = c.connect(); err==nil{
			break
		}
		l.Info("Rabbitmq trying to connect, attempts left: %d",i)
		time.Sleep(c.WaitTime)
	}
	if(err != nil){
		failOnError(err, "Couldn't connect to RabbitMQ")
	}
	return nil
}

func (c *Connection) connect() error {
    var err error

	c.Connection, err = amqp.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("amqp.Dial: %w", err)
	}

	c.Channel, err = c.Connection.Channel()
	if err != nil {
		return fmt.Errorf("c.Connection.Channel: %w", err)
	}

	err = c.Channel.ExchangeDeclare(
		c.ConsumerExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Connection.Channel: %w", err)
	}

	queue, err := c.Channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.QueueDeclare: %w", err)
	}

	err = c.Channel.QueueBind(
		queue.Name,
		"",
		c.ConsumerExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.QueueBind: %w", err)
	}

	c.Delivery, err = c.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("c.Channel.Consume: %w", err)
	}

	return nil
}
