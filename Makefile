.PHONY: test

test:
	go test

cert:
	go run cmd/gencert/*.go

server:
	go run cmd/server/*.go

image:
	docker build -t test-saml-idp .
