FROM golang:1.24 as build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app

ENV GIN_MODE=release
ENV PORT=8080
ENV BASE_URL=http://localhost:$PORT

COPY --from=build /src/server ./
COPY templates ./templates

EXPOSE $PORT

ENTRYPOINT ["./server"]
