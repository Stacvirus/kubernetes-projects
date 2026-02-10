
### Launching a k3s on docker using k3d

`k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2`

### Using PV on a k8s cluster
- Creating directory to hold actual data(locally k3d)
`docker exec k3d-k3s-default-agent-0 mkdir -p /tmp/kube`

### Database Setup in Kubernetes (Manual SQL Initialization)

This document explains how to manually initialize the PostgreSQL database running inside the Kubernetes cluster.
This method is useful during development or when you need to quickly bootstrap your tables without using migration tools or automation.

Note:
This is a quick and manual method. For production environments, automated migrations (e.g., Kubernetes Jobs or a migration CI pipeline) are recommended.

## 🚀 Overview

- We have deployed PostgreSQL into the Kubernetes cluster using a StatefulSet:
```
Namespace: exercises

StatefulSet: postgres-stset

Service: postgres-svc

Database name: pingpong

User: stac

Password: password
```

- PostgreSQL automatically creates the user and database on first startup based on the environment variables:
```
env:
  - name: POSTGRES_USER
    value: "stac"
  - name: POSTGRES_PASSWORD
    value: "password"
  - name: POSTGRES_DB
    value: "pingpong"
```

However, tables and schema are NOT created automatically.
We must create them manually using SQL.

# 🧩 Step 1 — Exec Into the PostgreSQL Pod

- Find the pod:

`kubectl get pods -n exercises`


- It should look like:

`postgres-stset-0`


- Open an interactive shell inside the database pod:

`kubectl exec -it -n exercises postgres-stset-0 -- sh`

# 🧩 Step 2 — Connect to PostgreSQL

- Once inside the pod, connect to the pingpong database using the default user:

`psql -U stac pingpong`


- If the connection succeeds, you will see:

`pingpong=#`

# 🧩 Step 3 — Create Required Tables

- Run the SQL statements directly inside the psql console.

Create the lines table:
```
CREATE TABLE IF NOT EXISTS lines (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

Confirm the table exists:

`\d lines`


- You should see the schema.

🧪 Step 4 — Test the Application

You can now run your application.
The table exists, so queries like INSERT or SELECT should work without errors.

Check logs:

`kubectl logs -n exercises <your-app-pod>`


You should no longer see:

pq: relation "lines" does not exist

🧹 Optional — Reset the Database

If you want to drop the table during development:

`DROP TABLE lines;`


Or drop and recreate the whole database (inside postgres connection):
```
psql -U stac
DROP DATABASE pingpong;
CREATE DATABASE pingpong;
```

## Intalling prometheus stack using helm

```
$ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
$ helm repo add stable https://charts.helm.sh/stable
$ helm repo update
```

```
$ kubectl create namespace prometheus
$ helm install prometheus-community/kube-prometheus-stack --generate-name --namespace prometheus
```
`kubectl -n prometheus port-forward prometheus-kube-prometheus-stack-1766-prometheus-0 9090:9090`

# Adding loki for logs:
```
$ helm repo add grafana https://grafana.github.io/helm-charts
$ helm repo update
$ kubectl create namespace loki-stack
  namespace/loki-stack created

$ helm upgrade --install loki --namespace=loki-stack grafana/loki-stack --set loki.image.tag=2.9.3

$ kubectl get all -n loki-stack
```
# Example usage: Access grafana UI
`$ kubectl -n prometheus port-forward kube-prometheus-stack-1602180058-grafana-59cd48d794-4459m 3000`

# Get Grafana 'admin' user password by running:

`  kubectl --namespace prometheus get secrets kube-prometheus-stack-1766116194-grafana -o jsonpath="{.data.admin-password}" | base64 -d ; echo `

### Install argo rollout for custom deployment strategies in k8s clusters(CRDs)
```
$ kubectl create namespace argo-rollouts
$ kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
```
- Install Argo kubectl cli tool
```
$ curl -LO https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-amd64
$ chmod +x kubectl-argo-rollouts-linux-amd64
$ sudo mv kubectl-argo-rollouts-linux-amd64 /usr/local/bin/kubectl-argo-rollouts
`