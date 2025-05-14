package server

import (
	"post_service/pkg/logger"
	"time"
)

type Option func(*Server)

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func Logger(l logger.Interface) Option {
	return func(s *Server) {
		s.logger = l
	}
}
