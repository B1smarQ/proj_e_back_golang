package server

import (
	"post_service/pkg/logger"
	"post_service/pkg/rabbitmq"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CallHandler func(*amqp.Delivery) (interface{}, error)
type Server struct {
	conn  *rabbitmq.Connection
	error chan error
	stop  chan struct{}

	router map[string]CallHandler

	timeout time.Duration

	logger logger.Interface
}

func New(url, serverExchange string, router []CallHandler, l logger.Interface, opts ...Option) *Server {
	server := &Server{
		conn:    nil,
		error:   make(chan error),
		stop:    make(chan struct{}),
		router:  make(map[string]CallHandler),
		timeout: 10 * time.Second,
		logger:  l,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}
