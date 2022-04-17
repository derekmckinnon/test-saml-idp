package main

import "github.com/crewjam/saml"

type ServiceProvider struct {
	Id               string
	EntityDescriptor saml.EntityDescriptor
}
