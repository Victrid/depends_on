# depends_on for kubernetes

This init container uses annotations to realize docker-compose style `depends_on`.

## Usage

### RBAC
Create ClusterRole to your cluster to allow the init container to check the status of the resources
Check [example/clusterrole.yaml](example/clusterrole.yaml).

```bash
kubectl apply -f example/clusterrole.yaml
```

You can authorize the default service account by binding such ClusterRole to the default service account in the namespace
where you want to use this init container. Check [example/clusterrolebinding.yaml](example/clusterrolebinding.yaml). This
is not recommended for production use.

Otherwise, create ClusterRoleBinding to bind the ClusterRole to the service account of the pod. Replace the `<namespace>` with the
namespace of your pod. For each namespace you'll need a dedicated ServiceAccount and bind it to the ClusterRole. Check
[example/serviceaccount.yaml](example/serviceaccount.yaml).

```bash
kubectl apply -f example/serviceaccount.yaml
```

This will allow not only the init container but following container to access the resources as ServiceAccount is a pod-level
resource. You can mount a separate ServiceAccount token to the init container if you want. 
Check [kubernetes/issues/66020](https://github.com/kubernetes/kubernetes/issues/66020) for more details.

### Pod Configuration

Use annotations as below on the pod:

```yaml
metadata:
  annotations:
    victrid.dev/depends_on: "<depends_on>"
```

The dependency can be configured using plain text or ConfigMap if too long:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dependency_config_map
data:
  depends_on: |
    - resource: "<resource>"
      name: "<name>"
      namespace: "<namespace>" # optional
      status: "<status>"       # optional
      raw: "<dependency>"      # or specify raw dependency
    - ...
```

Format:

```
depends_on        ::== depends_on_string | "ConfigMap=" locator
depends_on_string ::== dependency "," depends_on_string | dependency

dependency        ::== resource ":" locator | resource ":" locator "?" status
locator           ::== namespace "/" name | name
```

### Available resource: status pair

- CoreV1
  - "pod": "running", "ready"
  - "service": "ready", "running", "allready", "allrunning"
  
    Service selects pods by its selector, and inspect the status of the pods.
- AppsV1
  - "deployment", "statefulset", "replicaset": "ready", "available", "allready", "allavailable"
  - "daemonset": "ready", "available"


## Example

```yaml
# ConfigMap optional if having long dependencies
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
  namespace: example
data:
  depends_on: |
    - raw: "service:postgresql?available"
    - resource: "deployment"
      name: "prometheus"
      namespace: "metric"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: application
spec:
  selector:
    matchLabels:
      app: application
  template:
    metadata:
      name: application
      labels:
        app: application
      annotations:
        victrid.dev/depends_on: "service:postgresql?available,deployment:metric/prometheus"
        # or
        # victrid.dev/depends_on: "ConfigMap=example/config"
    spec:
      initContainers:
        - name: dependency
          image: victrid/depends_on
      # Remember to bind the service account if you created a dedicated one
      serviceAccountName: dependson-sa
      containers:
        - name: application
          image: some/application
```
