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
go run ./cmd/server
```

This will launch the IdP on port `8080` by default.

# Docker Image

Each tagged release is published on [GHCR](https://github.com/derekmckinnon/test-saml-idp/pkgs/container/test-saml-idp).

Simply `docker pull` or add the image to your `docker-compose.yml` file.

To run the container, you will need to volume mount your customized `config.yml` file into `/app/config.yml`.

Future releases should add a lot more configuration options (`env`s and flags) to enable more convenient configuration
scenarios.

If you have a specific architecture in mind that isn't currently supported, please
[open an Issue](https://github.com/derekmckinnon/test-saml-idp/issues/new). PRs are welcomed too :upside_down_face:
