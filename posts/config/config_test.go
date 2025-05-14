package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	// Test MySQL configuration
	if cfg.MySql.User != "root:B1smarQ._@/project_e" {
		t.Errorf("expected MySQL user %s, got %s", "root:B1smarQ._@/project_e", cfg.MySql.User)
	}

	// Test MongoDB configuration
	if cfg.MongoDB.URI != "mongodb://localhost:27017/project_e" {
		t.Errorf("expected MongoDB URI %s, got %s", "mongodb://localhost:27017/project_e", cfg.MongoDB.URI)
	}

	// Test Redis configuration
	if cfg.Redis.URI != "redis://localhost:6379" {
		t.Errorf("expected Redis URI %s, got %s", "redis://localhost:6379", cfg.Redis.URI)
	}

	// Test RabbitMQ configuration
	if cfg.RabbitMQ.URI != "amqp://localhost:5672" {
		t.Errorf("expected RabbitMQ URI %s, got %s", "amqp://localhost:5672", cfg.RabbitMQ.URI)
	}

	// Test App configuration
	if cfg.App.Port != "8080" {
		t.Errorf("expected App port %s, got %s", "8080", cfg.App.Port)
	}
	if cfg.App.Env != "development" {
		t.Errorf("expected App env %s, got %s", "development", cfg.App.Env)
	}
	if cfg.App.Host != "localhost" {
		t.Errorf("expected App host %s, got %s", "localhost", cfg.App.Host)
	}
}
