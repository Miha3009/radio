FROM golang:1.19 as builder
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /usr/src/app

COPY . .
RUN go mod download

WORKDIR /usr/src/app/cmd
RUN go build -o /usr/local/bin/radio -buildvcs=false

FROM alpine:latest
COPY --from=builder /usr/local/bin/radio ./
COPY --from=builder /usr/src/app/config.yaml ./config.yaml
CMD ["./radio"]
EXPOSE 8080
