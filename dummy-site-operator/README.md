# DummySite Operator

A Kubernetes operator that fetches a webpage from a URL and serves it via an nginx Deployment. Create a `DummySite` custom resource with a `websiteURL`, and the operator creates a ConfigMap (with the HTML content), Deployment, and Service in that namespace.

---

## Table of Contents

- [Custom Resource Definition (CRD)](#custom-resource-definition-crd)
- [RBAC: ServiceAccount, ClusterRole, ClusterRoleBinding](#rbac-serviceaccount-clusterrole-clusterrolebinding)
- [Package Usage](#package-usage)
- [Code Generation](#code-generation)
- [Deployment](#deployment)

---

## Custom Resource Definition (CRD)

### What is a CRD?

A **Custom Resource Definition** extends the Kubernetes API by defining new resource types. Instead of only using built-in resources (Pods, Services, Deployments, etc.), you can define your own—in this case, `DummySite`.

### DummySite CRD Structure

| Field | Description |
|-------|-------------|
| **Group** | `dummy-site-operator.stac.dev` |
| **Version** | `v1alpha1` |
| **Kind** | `DummySite` |
| **Plural** | `dummysites` |
| **Scope** | `Namespaced` (DummySite instances exist within a namespace) |

### Spec and Status

- **Spec** (`spec`): User-provided desired state
  - `website_url`: URL to fetch and serve (e.g., `https://example.com`)
- **Status** (`status`): Operator-maintained actual state
  - `ready`: Whether the site is ready to serve

### Subresources

The CRD defines a **status subresource**. This lets the operator update `status` independently from `spec`, and prevents `status` changes from triggering unnecessary reconciles.

---

## RBAC: ServiceAccount, ClusterRole, ClusterRoleBinding

### Why RBAC?

The controller runs as a Pod and needs permission to talk to the Kubernetes API (list DummySites, create Deployments, etc.). RBAC (Role-Based Access Control) defines what the controller is allowed to do.

### ServiceAccount

A **ServiceAccount** is an identity for processes running inside a Pod. When the controller Pod starts, it uses the `dummysite-controller` ServiceAccount. Every API call from the Pod is authenticated as that ServiceAccount.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: dummysite-controller
  namespace: default
```

The controller Deployment sets `serviceAccountName: dummysite-controller`, so the Pod runs under this identity.

### ClusterRole vs Role

| | Role | ClusterRole |
|---|-----|-------------|
| **Scope** | Namespace-scoped | Cluster-wide |
| **Grants access to** | Resources in a single namespace | Resources across all namespaces |

**Why ClusterRole here?**

The controller uses controller-runtime’s default manager, which uses a cluster-wide cache. It needs to:

- **List and watch** `DummySite` resources across all namespaces
- **Create** Pods, Services, ConfigMaps, and Deployments in the same namespace as each DummySite (which can be any namespace)

Those operations are cluster-scoped (list/watch) or can target any namespace (create). A namespace-scoped `Role` only grants permissions in one namespace (e.g. `default`), so the controller would get “forbidden” when listing at cluster scope. A **ClusterRole** grants the needed permissions across the cluster.

### ClusterRoleBinding

A **ClusterRoleBinding** connects a subject (here, the ServiceAccount) to a ClusterRole. It grants the ClusterRole’s permissions to that subject.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: dummysite-controller
subjects:
- kind: ServiceAccount
  name: dummysite-controller
  namespace: default
roleRef:
  kind: ClusterRole
  name: dummysite-controller
  apiGroup: rbac.authorization.k8s.io
```

The subject must include `namespace` because ServiceAccounts are namespace-scoped. `ClusterRoleBinding` itself is cluster-scoped and has no `metadata.namespace`.

---

## Package Usage

| Package | Purpose |
|---------|---------|
| `github.com/stacvirus/dummy-site-operator/api/v1alpha1` | DummySite CRD Go types, scheme registration, and DeepCopy implementations |
| `github.com/stacvirus/dummy-site-operator/controllers` | DummySite reconciler logic |
| `k8s.io/api` | Core Kubernetes API types (Deployment, Service, ConfigMap, Pod, etc.) |
| `k8s.io/apimachinery` | `metav1`, `runtime`, `intstr`, schema, etc. |
| `sigs.k8s.io/controller-runtime` | Manager, controller, client, logging, and reconciliation framework |

### Flow

1. **main.go** – Creates the manager, registers the DummySite scheme, wires the reconciler, starts the manager.
2. **api/v1alpha1** – Defines the DummySite CRD types and registers them with the runtime scheme.
3. **controllers** – Implements `Reconcile()`: fetches the URL, creates ConfigMap, Deployment, and Service.

---

## Code Generation

The CRD manifest and DeepCopy code are **generated** from the Go types and markers. Do not edit them manually.

### Tools

- **controller-gen** (from [controller-tools](https://github.com/kubernetes-sigs/controller-tools)) reads markers like `// +kubebuilder:object:root=true` and generates CRDs and DeepCopy methods.

### Install controller-gen

```bash
go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest
```

Ensure `$GOPATH/bin` (or `$HOME/go/bin`) is in your `PATH`.

### Generate `zz_generated.deepcopy.go`

DeepCopy implementations are required for types used with the controller-runtime cache and client. The `// +kubebuilder:object:root=true` marker tells controller-gen to generate them.

```bash
controller-gen object paths="./api/..."
```

This creates/updates `api/v1alpha1/zz_generated.deepcopy.go` with `DeepCopy()`, `DeepCopyInto()`, and `DeepCopyObject()` for `DummySite` and `DummySiteList`.

### Generate `manifests/crd/_dummysites.yaml`

The CRD YAML is generated from the same Go types and markers:

```bash
controller-gen crd paths="./api/..." output:crd:dir=./manifests/crd
```

This produces `manifests/crd/dummy-site-operator.stac.dev_dummysites.yaml` (or similar). The `_dummysites.yaml` file in this repo is the generated CRD, possibly renamed or organized for deployment.

### Generate both in one command

```bash
controller-gen object crd paths="./api/..." output:crd:dir=./manifests/crd
```

### Kubebuilder markers used

| Marker | Effect |
|--------|--------|
| `// +kubebuilder:object:root=true` | Marks as a root type; generates CRD and DeepCopy |
| `// +kubebuilder:subresource:status` | Adds a status subresource to the CRD |

---

## Deployment

### 1. Apply the CRD

```bash
kubectl apply -f manifests/crd/
```

### 2. Apply RBAC and controller

```bash
kubectl apply -f manifests/rbac.yaml
kubectl apply -f manifests/controller.yaml
```

### 3. Create a DummySite

```bash
kubectl apply -f manifests/dummysite.yaml
```

### 4. Port-forward and access

```bash
kubectl port-forward svc/example 8080:80
```

Then open `http://localhost:8080`.

---

## JSON Tag and CRD Consistency

The Go struct’s `json` tag must match the field name in the CRD schema and in manifests. controller-gen derives the CRD schema from the struct’s `json` tags.

For example, if the struct uses `json:"website_url"`, the CRD and manifests should use `website_url`. If you change the `json` tag, regenerate the CRD and update any existing manifests accordingly.
