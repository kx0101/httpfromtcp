following:

- [RFC-9112](https://datatracker.ietf.org/doc/html/rfc9112)
- [RFC-9110](https://datatracker.ietf.org/doc/html/rfc9110)

open two terminals:

1. `go run ./cmd/tcplistener | tee /tmp/headers.txt`
2. `curl -X POST http://localhost:42069/coffee \
-H 'Content-Type: application/json' \
-d '{"type": "dark mode", "size": "medium"}'`

check the headers.txt file, you will see the following:

```
Request Line:
  - Method: POST
  - Target: /coffee
  - Version: HTTP/1.1
Headers:
  - host: localhost:42069
  - user-agent: curl/8.5.0
  - accept: */*
  - content-type: application/json
  - content-length: 39
Body:
{"type": "dark mode", "size": "medium"}
```
