package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Service struct {
	Id string `yaml:"id"`
}

type User struct {
	Username  string `yaml:"username"`
	Email     string `yaml:"email"`
	Password  string `yaml:"password"`
	FirstName string `yaml:"first_name"`
	LastName  string `yaml:"last_name"`
}

type Config struct {
	Services []Service
	Users    []User
}

func (config *Config) InitFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}
