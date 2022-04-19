Test SAML IdP
=============

A simple and configurable SAML IdP for testing and development scenarios.

The IdP is built using:

- [crewjam/saml](https://github.com/crewjam/saml)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)

**Do not use this in production!**

# Getting Started

This project requires Go 1.18+ for development.

1. `cp config.example.yml config.yml`
2. Populate `config.yml` with your own service provider and user configuration

# Usage

To run locally:

```shell
go run main.go
```

This will launch the IdP on port `8080` by default.
