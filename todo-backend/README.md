### Using NATS for our EDA go server

## Install go nats library
`go get github.com/nats-io/nats.go`

## Get Nats credentials from k8s
`kubectl get secret my-nats -o yaml`