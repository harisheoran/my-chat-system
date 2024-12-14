## my-chat-system
Scaleable personal chat system

## Connect to my instance
Visit https://chat.harisheoran.xyz/

OR
```
wscat -c https://chat.harisheoran.xyz/
```

Test Deployments
- http://13.203.105.149:1317/v1/home

### How to set up?
- Install [***wscat***](https://github.com/websockets/wscat) tool

- Clone the repo and go to *chat-server* directory

``` go mod download ```

``` go run main.go ```

- Run the command - ``` wscat -c ws://localhost:1316/v1/chat ```

- Start chatting
