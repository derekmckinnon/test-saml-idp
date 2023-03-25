FROM golang:1.20 as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./cmd/server

FROM gcr.io/distroless/base

WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080
ENV BASE_URL=http://localhost:$PORT

COPY --from=build /src/server ./
COPY *.pem ./
COPY templates ./templates

EXPOSE $PORT

ENTRYPOINT ["./server"]
