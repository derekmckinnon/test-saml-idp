package idp

import "fmt"

type Config struct {
	Host     string
	Port     string
	Services []Service
	Users    []User
}

func (c Config) BaseUrl() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

type Service struct {
	EntityId                 string `mapstructure:"entity_id"`
	AssertionConsumerService string `mapstructure:"assertion_consumer_service"`
}

type User struct {
	Username  string
	Email     string
	Password  string
	FirstName string `mapstructure:"first_name"`
	LastName  string `mapstructure:"last_name"`
}
