package config

type (
	Config struct {
		MySql    MySql
		MongoDB  MongoDB
		Redis    Redis
		RabbitMQ RabbitMQ
		App      App
	}

	MySql struct {
		User string `default:"root:B1smarQ._@/project_e"`
	}

	MongoDB struct {
		URI string `default:"mongodb://localhost:27017/project_e"`
	}

	Redis struct {
		URI string `default:"redis://localhost:6379"`
	}

	RabbitMQ struct {
		URI string `default:"amqp://guest:guest@localhost:5672"`
	}

	App struct {
		Port     string `default:"8080"`
		Env      string `default:"development"`
		Host     string `default:"localhost"`
		Name     string `default:"project_e"`
		Version  string `default:"1.0.0"`
		LogLevel string `default:"debug"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	cfg.MySql.User = "root:B1smarQ._@/project_e"
	cfg.MongoDB.URI = "mongodb://localhost:27017/project_e"
	cfg.Redis.URI = "redis://localhost:6379"
	cfg.RabbitMQ.URI = "amqp://localhost:5672"
	cfg.App.Port = "8080"
	cfg.App.Env = "development"
	cfg.App.Host = "localhost"

	return cfg
}
