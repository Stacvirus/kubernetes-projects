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
`kubectl get secret my-nats -o jsonpath='{.data.nats-server\.conf}' | base64 --decode`

### Installing ArgoCD in our k8s cluster
- create a name space
`$ kubectl create namespace argocd`

- Apply argocd manifest files
`$ kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml`

- Openning access to Argocd
`$ kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'`
or locally:
`$ kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'`
simplest: `$ kubectl port-forward svc/argocd-server -n argocd 8080:443`

- Getting credentials(decode the result as it a base64 value)
`$ kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 --decode`

### GitOps Installing ArgoCD in our k8s cluster
`
$ kubectl create namespace argocd
$ kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
`
- Open access via loadBalancer
`$ kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'`

- If you are on a local cluster(k3d)
`kubectl port-forward svc/argocd-server -n argocd 8080:443`

- get the admin password
`kubectl get -n argocd secrets argocd-initial-admin-secret -o yaml`

- get the actual password string
`kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d && echo`