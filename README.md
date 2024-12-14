## my-chat-system
Scaleable personal chat system

## Setup locally
- ``` cd chat-server ```
- ``` go mod download ```
- Install [air](https://github.com/air-verse/air) for fast reload
- start the server
```
air --build.cmd "go build -o bin/api ./cmd/server" --build.bin "./bin/api -port=3000 -env=test
```

## Connect to my instance
```
wscat -c http://13.203.105.149:1316/v1/chat
```

### How to set up?
- Install [***wscat***](https://github.com/websockets/wscat) tool

- Clone the repo and go to *chat-server* directory

``` go mod download ```

``` go run main.go ```

- Run the command - ``` wscat -c ws://localhost:1316/v1/chat ```

- Start chatting
