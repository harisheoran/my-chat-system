FROM golang:1.23.3 as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/my-chat-system-server .

FROM alpine:3.20 as prod

WORKDIR /app

COPY --from=build /app/bin/my-chat-system-server /app/bin/my-chat-system-server

RUN chmod +x /app/bin/my-chat-system-server

EXPOSE 1316

ENTRYPOINT [ "/app/bin/my-chat-system-server" ]
