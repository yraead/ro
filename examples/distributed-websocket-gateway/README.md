
# Distributed websocket gateway

Websocket client <-> Websocket gateway <-> Redis

Can be used for chats.

## Run

```bash
go run *.go
```

Then, open multiple websocket clients:

```bash
wscat -c 'ws://localhost:8080?room=1'
wscat -c 'ws://localhost:8080?room=1'
wscat -c 'ws://localhost:8080?room=abcd'
wscat -c 'ws://localhost:8080?room=1'
```
