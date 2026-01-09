# Todo App

A simple Todo application written in Go with configurable port settings through environment variables.

## Prerequisites

- Go 1.24.5 or higher
- Docker (optional)

## Project Structure

```
.
├── Dockerfile
├── go.mod
├── main.go
├── .env
└── README.md
```

## Environment Variables

Create a `.env` file in the project root:

```plaintext
PORT=8000
```

## Running Locally

1. Clone the repository:
```bash
git clone https://github.com:Stacvirus/hash-generator-app/tree/1.2
cd todo-app
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on the port specified in your `.env` file (default: 8000)

## Running with Docker

1. Build the Docker image:
```bash
docker build -t todo-app .
```

2. Run the container:
```bash
# Using default port (8000)
docker run -p 8000:8000 todo-app

# Override port using environment variable
docker run -p 3000:3000 -e PORT=3000 todo-app

# Using custom .env file
docker run -p 8000:8000 -v $(pwd)/.env:/app/.env todo-app
```

## API Endpoints

- `GET /`: Returns a simple message indicating the app is running

## Development

### Modifying the Port

You can change the port in three ways:

1. Update the `.env` file
2. Set the PORT environment variable:
```bash
PORT=3000 go run main.go
```
3. Use Docker environment variables as shown above

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port number | 8000 |

## Running app container on k8s pod using kustomize
`kubectl apply -k .`

## Installing NATS on k8s cluster
`$ helm install --set metrics.enabled=true,auth.enabled=false my-nats oci://registry-1.docker.io/bitnamicharts/nats`

- here is the output logs you obtain if everything works properly:
```
Pulled: registry-1.docker.io/bitnamicharts/nats:9.0.28
Digest: sha256:9ac9ba1073909fa55ee47e79520cb3c4abd56be20a7bc5e6dad24bca15c732e9
NAME: my-nats
LAST DEPLOYED: Tue Jan  6 04:21:08 2026
NAMESPACE: default
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: nats
CHART VERSION: 9.0.28
APP VERSION: 2.11.8

⚠ WARNING: Since August 28th, 2025, only a limited subset of images/charts are available for free.
    Subscribe to Bitnami Secure Images to receive continued support and security updates.
    More info at https://bitnami.com and https://github.com/bitnami/containers/issues/83267

** Please be patient while the chart is being deployed **

NATS can be accessed via port 4222 on the following DNS name from within your cluster:

   my-nats.default.svc.cluster.local

NATS monitoring service can be accessed via port 8222 on the following DNS name from within your cluster:

    my-nats.default.svc.cluster.local

You can create a pod to be used as a NATS client:

    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: my-nats-client
  namespace: default
spec:
  containers:
  - name: cli
    image: docker.io/bitnami/natscli
    command: ["sleep", "infinity"]
    env:
    - name: NATS_USER
      value: "$NATS_USER"
    - name: NATS_PASS
      value: "$NATS_PASS"
EOF

Then, access the pod and connect to NATS:

    kubectl exec --tty -i my-nats-client --namespace default -- bash
    nats -s nats://my-nats.default.svc.cluster.local:4222  subscribe SomeSubject
    nats -s nats://my-nats.default.svc.cluster.local:4222  publish SomeSubject "Some message"

To access the Monitoring svc from outside the cluster, follow the steps below:

1. Get the NATS monitoring URL by running:

    echo "Monitoring URL: http://127.0.0.1:8222"
    kubectl port-forward --namespace default svc/my-nats 8222:8222

2. Open a browser and access the NATS monitoring browsing to the Monitoring URL

3. Get the NATS Prometheus Metrics URL by running:

    echo "Prometheus Metrics URL: http://127.0.0.1:7777/metrics"
    kubectl port-forward --namespace default svc/my-nats-metrics 7777:7777

4. Access NATS Prometheus metrics by opening the URL obtained in a browser.

WARNING: There are "resources" sections in the chart not set. Using "resourcesPreset" is not recommended for production. For production installations, please set the following values according to your workload needs:
  - metrics.resources
  - resources
+info https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
```

- An altenative to keep a more maintainable helm will be create a helm values file(natValues.yml)
```
auth:
  enabled: <strong>false</strong>

metrics:
  enabled: <strong>true</strong>
  serviceMonitor:
    enabled: <strong>true</strong>
    namespace: prometheus
```
now just use the command: `helm upgrade -f helm/natvalues.yml my-nats oci://registry-1.docker.io/bitnamicharts/nats`
NB: you could also overrides some other values like the image names and tags