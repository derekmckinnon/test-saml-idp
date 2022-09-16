package idp

type Config struct {
	Host     string
	Services []Service
	Users    []User
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
