FROM golang:1.18 as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./test-saml-idp .

FROM gcr.io/distroless/base

WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080
ENV BASE_URL=http://localhost:$PORT

COPY --from=build /src/test-saml-idp ./
COPY *.pem ./
COPY templates ./templates

EXPOSE $PORT

ENTRYPOINT ["./test-saml-idp"]
