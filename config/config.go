package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	RabbitMQ RabbitMQ       `yaml:"rabbitmq"`
	App      App            `yaml:"app"`
	Metrics  Metrics        `yaml:"metrics"`
}

type PostgresConfig struct {
	Host     string `yaml:"PostgresqlHost" env:"POSTGRESQL_HOST"`
	Port     string `yaml:"PostgresqlPort" env:"POSTGRESQL_PORT"`
	User     string `yaml:"PostgresqlUser" env:"POSTGRESQL_USERNAME"`
	Password string `yaml:"PostgresqlPassword" env:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"PostgresqlDbname" env:"POSTGRESQL_NAME"`
}

type RabbitMQ struct {
	Username     string `yaml:"username" env-required:"true"`
	Password     string `yaml:"password" env-required:"true"`
	Host         string `yaml:"host" env-required:"true"`
	Port         string `yaml:"port" env-required:"true"`
	ExchangeName string `yaml:"exchangeName" env-required:"true"`
	QueueName    string `yaml:"queueName" env-required:"true"`
	ConsumerTag  string `yaml:"consumerTag" env-required:"true"`
	BindingKey   string `yaml:"bindingKey" env-required:"true"`
}

type Metrics struct {
	Jaeger Jaeger `yaml:"jaeger"`
}

type Jaeger struct {
	Endpoint string `yaml:"endpoint"`
}

type App struct {
	Port string `yaml:"port"`
}

func LoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}
	return &cfg

}
