# 📘 Istio + Bookinfo Sample with Gateway API — Deep Dive README

This document explains how to install Istio, deploy the Bookinfo sample app on a local Kubernetes cluster (using k3d), expose it via the Kubernetes Gateway API, and exercise it with basic traffic. It covers everything you did step-by-step, including traps and solutions.

## 📌 1. Download & Install Istio

```bash
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.28.3
export PATH=$PWD/bin:$PATH
```

**Explanation:**

- This uses the official Istio installer script to download the specified version (1.28.3).
- It places the `istioctl` binary in `istio-1.28.3/bin`.
- You export that path so you can run `istioctl` directly.

```bash
istioctl version
```

At first this failed because no cluster was running.

`istioctl` tries to contact the Kubernetes API to check for control plane components. If there is no cluster or no Istio pods, you'll get connection refused or "Istio not present" messages.

## 📌 2. Start Kubernetes Local Cluster (k3d)

```bash
k3d cluster create \
  --port 8082:30080@agent:0 \
  -p 8081:80@loadbalancer \
  --agents 2 \
  --k3s-arg '--disable=traefik@server:*'
```

**Explanation:**

- Creates a k3d (k3s) cluster with 2 agents.
- Maps:
  - Agent port 30080 → local 8082
  - LoadBalancer port 80 → local 8081
- Disables the built-in Traefik ingress (so Istio can handle ingress).

This gives you a clean cluster suitable for Istio. There are no Istio pods yet.

## 📌 3. Install Istio into the Cluster

```bash
istioctl install --set profile=ambient --set values.global.platform=k3d -y
```

**Explanation:**

- Uses the "ambient" profile — this is Istio's sidecar-less mode.
- In ambient mode, workloads don't run with sidecar proxies; instead ambient proxies (Ztunneld) handle traffic.
- We explicitly set platform to k3d to help Istio configure itself correctly.

**Key Point:** The official docs show a "demo" profile for newcomers, but ambient is great for sidecar-less experimentation.

## 📌 4. Label Namespace for Sidecar Injection

```bash
kubectl label namespace default istio-injection=enabled
```

**Why?**

This tells Istio that new pods in `default` should have sidecar proxies injected automatically (into the pod spec). It's necessary when using sidecar mode.

For ambient mode it's less critical but still good habit.

## 📌 5. Deploy Bookinfo Sample App

```bash
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
kubectl apply -f samples/bookinfo/platform/kube/bookinfo-versions.yaml
```

**What this does:**

- Deploys all Bookinfo microservices (productpage, details, ratings, reviews v1/v2/v3).
- After a short while, the pods should show `READY 2/2 RUNNING`.

This matches the official example.

## 📌 6. Verify Bookinfo from Inside Cluster

```bash
kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings \
  -- curl -sS productpage:9080/productpage | grep "<title>"
```

**Explanation:**

Executes a curl inside the Kubernetes network, bypassing any ingress. This verifies that Bookinfo is functioning internally.

## 📌 7. Expose Bookinfo with Kubernetes Gateway API

### 7A. Install Gateway API CRDs

```bash
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/latest/download/standard-install.yaml
```

**Why this is necessary:**

The Gateway API resources like `Gateway` and `HTTPRoute` are not built-in Kubernetes objects — they're CRDs maintained by SIG Networking. Most clusters don't have them by default.

### 7B. Apply the Gateway API Resources

```bash
kubectl apply -f samples/bookinfo/gateway-api/bookinfo-gateway.yaml
kubectl annotate gateway bookinfo-gateway networking.istio.io/service-type=ClusterIP
```

**Explanation:**

- Installs a Kubernetes Gateway ("bookinfo-gateway").
- Installs an HTTPRoute that ties the gateway to the Bookinfo services.
- Annotating the gateway sets its service type to ClusterIP, so we can port-forward to it instead of using a LoadBalancer.

This is the same pattern shown in the official docs.

## 📌 8. Test Access Through Port-Forward

```bash
kubectl port-forward svc/bookinfo-gateway-istio 8080:80
```

Now point a browser or curl at:

```
http://localhost:8080/productpage
```

You should see Bookinfo load with changing reviews — this is because traffic is routed across review versions.

## 📌 9. Install Monitoring Addons

```bash
kubectl apply -f samples/addons/prometheus.yaml
kubectl apply -f samples/addons/kiali.yaml
```

Then start Kiali:

```bash
istioctl dashboard kiali
```

This should open:

```
http://localhost:20001/kiali
```

Addons provide telemetry and visualization of the mesh. This matches the "View the dashboard" section in the official guide.

## 📌 10. Generating Traffic

You tried:

```bash
for i in $(seq 1 100); do curl -sSI -o /dev/null http://localhost:8080/productpage; done
```

Sometimes that fails if the port-forward isn't active or gateway not reachable.

You also ran:

```bash
oha -z 10s -q 100 http://localhost:8080/productpage
```

**Explanation:**

- Sends ~100 requests/sec for 10s (≈1000 requests total).
- Great for load testing and ensuring gateways handle throughput correctly.

## 📌 11. Authorization Policies in Ambient Mode

Authorization policies in Istio ambient mode work differently than in sidecar mode. In ambient mode, there are **two layers** where you can apply policies:

1. **Ztunnel layer** (L4 - secure overlay) - handles encrypted traffic between pods
2. **Waypoint proxy layer** (L7 - optional) - handles advanced routing and policies

### 11A. Ztunnel Authorization Policy

First, you created a policy to allow the gateway to reach the productpage at the **ztunnel layer**:

```bash
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: productpage-ztunnel
  namespace: default
spec:
  selector:
    matchLabels:
      app: productpage
  action: ALLOW
  rules:
  - from:
    - source:
        principals:
        - cluster.local/ns/default/sa/bookinfo-gateway-istio
EOF
```

**What this does:**

- Uses `selector` to target pods with label `app: productpage`
- Allows traffic from the gateway's service account (`bookinfo-gateway-istio`)
- Enforces this at the **ztunnel (L4) layer** - the secure overlay network

**Why it's needed:**

In ambient mode, ztunnel creates an encrypted mesh between pods. By default, once you enable authorization policies, all traffic is denied unless explicitly allowed.

### 11B. Testing Access Denial

You deployed a test client to verify the policy:

```bash
kubectl apply -f samples/curl/curl.yaml
kubectl exec deploy/curl -- curl -s "http://productpage:9080/productpage"
```

**Result:** `RBAC: access denied`

This failed because the `curl` pod's service account wasn't in the allowed principals list.

### 11C. Waypoint Proxy for L7 Policies

For more advanced L7 (HTTP-level) policies, you deployed a **waypoint proxy**:

```bash
istioctl waypoint apply --enroll-namespace --wait
```

**What this does:**

- Creates a waypoint proxy (a special Envoy proxy) in the namespace
- Labels the namespace with `istio.io/use-waypoint: waypoint`
- All traffic in the namespace now flows through this waypoint for L7 processing

**Verify the waypoint:**

```bash
kubectl get gtw waypoint
```

Output shows the waypoint gateway is ready.

### 11D. Waypoint Authorization Policy

Then you created an L7 authorization policy attached to the waypoint:

```bash
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: productpage-waypoint
  namespace: default
spec:
  targetRefs:
  - kind: Service
    group: ""
    name: productpage
  action: ALLOW
  rules:
  - from:
    - source:
        principals:
        - cluster.local/ns/default/sa/curl
    to:
    - operation:
        methods: ["GET"]
EOF
```

**Key differences from ztunnel policy:**

- Uses `targetRefs` instead of `selector` - this targets a **Service** resource
- Can specify HTTP methods (`GET`) - this is L7-aware
- Allows the `curl` service account to make GET requests

**Important:** This policy is evaluated at the **waypoint proxy**, not at ztunnel.

### 11E. Updating Ztunnel Policy for Waypoint

Since traffic now flows through the waypoint, you had to update the ztunnel policy to allow the waypoint itself:

```bash
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1
kind: AuthorizationPolicy
metadata:
  name: productpage-ztunnel
  namespace: default
spec:
  selector:
    matchLabels:
      app: productpage
  action: ALLOW
  rules:
  - from:
    - source:
        principals:
        - cluster.local/ns/default/sa/bookinfo-gateway-istio
        - cluster.local/ns/default/sa/waypoint
EOF
```

**What changed:**

- Added `cluster.local/ns/default/sa/waypoint` to allowed principals
- Now the waypoint proxy can reach productpage at the ztunnel layer

**Traffic flow with waypoint:**

```
curl pod → waypoint proxy (L7 policy check) → ztunnel (L4 policy check) → productpage
```

After this, the curl pod could successfully access productpage:

```bash
kubectl exec deploy/curl -- curl -s http://productpage:9080/productpage | grep -o "<title>.*</title>"
```

### 11F. Traffic Splitting with HTTPRoute

You also created an HTTPRoute for traffic splitting between review versions:

```bash
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: reviews
spec:
  parentRefs:
  - group: ""
    kind: Service
    name: reviews
    port: 9080
  rules:
  - backendRefs:
    - name: reviews-v1
      port: 9080
      weight: 90
    - name: reviews-v2
      port: 9080
      weight: 10
EOF
```

**What this does:**

- Routes 90% of traffic to `reviews-v1`
- Routes 10% of traffic to `reviews-v2`
- This is processed at the **waypoint proxy** (L7 routing)

**Why HTTPRoute needs waypoint:**

Gateway API resources like HTTPRoute require L7 processing, which is only available at the waypoint proxy in ambient mode.

### 🔑 Key Concepts Summary

**Ztunnel (L4) Policies:**
- Use `selector` with pod labels
- Enforce identity-based access (which service account)
- Mandatory for securing the mesh

**Waypoint (L7) Policies:**
- Use `targetRefs` with Service resources
- Can inspect HTTP methods, headers, paths
- Required for advanced routing (traffic splitting, canary deployments)
- Optional - only deploy if you need L7 features

**Service Accounts and Principals:**
- Each pod runs with a Kubernetes service account
- Istio uses these as identities: `cluster.local/ns/<namespace>/sa/<service-account>`
- Policies match traffic based on source principal (identity)

## 📌 12. Cleanup

You deleted resources in reverse order:

```bash
kubectl delete -f samples/bookinfo/platform/kube/bookinfo.yaml
kubectl delete -f samples/bookinfo/gateway-api/bookinfo-gateway.yaml
```

This properly removes the application and gateways.

## 🧠 Key Concepts & Notes

### ✔ Istio Profiles

- **demo** — great for learning
- **ambient** — sidecar-less; good for experimentation without proxies injected into pods
- **production** — for real deployments

Official docs show the demo profile by default, but any profile works as long as gateways and proxies are present.

### ✔ Gateway API vs Istio APIs

Istio can use:

- ✔ **Kubernetes Gateway API** (modern, standard)
- ✔ **Istio APIs** (Gateway, VirtualService) (legacy, still widely used)

Both work, but Gateway API is the future.

## 📚 Further Reading

### 📘 Official Istio Docs

- **Bookinfo App Example** — https://istio.io/docs/examples/bookinfo/
- **Getting Started Guide** — https://istio.io/latest/docs/setup/getting-started/
- **Gateway API Details for Istio** — https://istio.io/latest/docs/ops/configuration/traffic-management/gateway/

### 🛠 Kubernetes Gateway API

- **Gateway API overview:** https://gateway.spec.networking.k8s.io/ — CRDs installed manually.