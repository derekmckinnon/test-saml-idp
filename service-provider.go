package main

import (
	"github.com/crewjam/saml"
	"net/http"
)

func (s *IdpServer) GetServiceProvider(_ *http.Request, id string) (*saml.EntityDescriptor, error) {
	service, err := s.Store.GetServiceProvider(id)
	if err != nil {
		return nil, err
	}

	return &service.Metadata, nil
}
