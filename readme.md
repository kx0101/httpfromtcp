open two terminals:

1. `go run ./cmd/tcplistener | tee /tmp/headers.txt`
2. `curl http://localhost:42069/use-neovim-btw`

check the headers.txt file, you will see the following:

```
Request Line:
  - Method: GET
  - Target: /use-neovim-btw
  - Version: HTTP/1.1
Headers:
  - accept: */*
  - host: localhost:42069
  - user-agent: curl/8.5.0
```
