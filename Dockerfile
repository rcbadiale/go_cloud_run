## Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cloud_run cmd/main.go

## Run stage
FROM alpine

# Settings
RUN apk add -U tzdata ca-certificates
ENV TZ=America/Sao_Paulo
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime

COPY --from=builder /app/cloud_run /cloud_run

ENTRYPOINT [ "/cloud_run" ]
