FROM golang:1.22-alpine3.19 as builder
WORKDIR /app
COPY . .
RUN go build

FROM alpine:3.19
ARG TIMEZONE=Asia/Taipei
WORKDIR /app
RUN apk add tzdata && \
    cp "/usr/share/zoneinfo/$TIMEZONE" /etc/localtime && \
    echo "$TIMEZONE" > /etc/timezone
RUN mkdir -p /app
COPY --from=builder /app/synochatgpt .
ENTRYPOINT ./synochatgpt
