FROM golang:1.12.5 as build

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64
RUN chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/harrisonlob/goyagi

COPY Gopkg.toml Gopkg.toml
COPY Gopkg.lock Gopkg.lock

RUN dep ensure -vendor-only

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o ./bin/goyagi ./cmd/serve/main.go

FROM alpine:3.8

COPY --from=build /go/src/github.com/harrisonlob/goyagi/bin .

CMD ["./goyagi"]
