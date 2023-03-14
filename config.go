package idp

type Config struct {
	Host             string           `mapstructure:"host"`
	Services         []Service        `mapstructure:"services"`
	Users            []User           `mapstructure:"users"`
	LoginPageOptions LoginPageOptions `mapstructure:"login_page"`
}

type Service struct {
	EntityId                 string `mapstructure:"entity_id"`
	AssertionConsumerService string `mapstructure:"assertion_consumer_service"`
}

type User struct {
	Username  string   `mapstructure:"username"`
	Email     string   `mapstructure:"email"`
	Password  string   `mapstructure:"password"`
	FirstName string   `mapstructure:"first_name"`
	LastName  string   `mapstructure:"last_name"`
	Groups    []string `mapstructure:"groups"`
}

type LoginPageOptions struct {
	Title       string `mapstructure:"title"`
	Description string `mapstructure:"description"`
	DumpUsers   bool   `mapstructure:"dump_users"`
}
