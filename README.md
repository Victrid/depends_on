# depends_on for kubernetes

This init container uses annotations to realize docker-compose style `depends_on`.

## Usage

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
resource          ::== "service" | "deployment" | "statefulset" | "daemonset" 
locator           ::== namespace "/" name | name
```

## Available status



## Example

```yaml
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
      containers:
      - name: application
        image: some/application
```
