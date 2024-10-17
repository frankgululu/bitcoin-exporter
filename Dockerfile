FROM golang:latest as builder
COPY . /src
RUN apt-get update && \
    cd /src && \
    GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/main main.go


FROM alpine:latest
WORKDIR "/src"
COPY ./config.yaml /src
COPY --from=builder /src/build/main /usr/bin/main
EXPOSE  2024
CMD ["/usr/bin/main"]