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


## File Structure
### remote
contains config files and setup scripts for prod.

### Makefile
will contain recipes for automating common administrative tasks — like
auditing our Go code, building binaries, and executing database migrations.

### internal
The internal directory will contain various ancillary packages used by our API. It will
contain the code for interacting with our database, doing data validation, sending emails
and so on. Basically, any code which isn’t application-specific and can potentially be
reused will live in here. Our Go code under cmd/api will import the packages in the
internal directory (but never the other way around).

### cmd/api
The cmd/api directory will contain the application-specific code.
This will include the code for running the server, reading and writing HTTP
requests, and managing authentication.
