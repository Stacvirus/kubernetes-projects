### Todo Backend usage
## Http Request
```
curl -X POST http://localhost:8080/todos   -H "Content-Type: application/json"   -d '{"task": "broadcast telegram bot test"}'
{"id":7,"task":"broadcast telegram bot test","done":false,"created_at":"2026-01-08T04:29:12.399685Z","updated_at":"2026-01-08T04:29:12.399685Z"}
```

### Using NATS for our EDA go server

## Install go nats library
`go get github.com/nats-io/nats.go`

## Get Nats credentials from k8s
`kubectl get secret my-nats -o yaml`