# Scenarios

KubeWise scenarios are hypothetical mutations applied to a snapshot of your cluster. Each scenario produces a cost estimate and risk assessment without making any actual changes.

## Scenario file format

All scenario files follow this structure:

```yaml
apiVersion: kubewise.io/v1alpha1
kind: RightSize | Consolidate | SpotMigrate | Composite
metadata:
  name: my-scenario
  description: "Human-readable description"
spec:
  # kind-specific fields
```

## RightSize

Adjusts resource requests and limits based on actual usage percentiles.

### YAML schema

```yaml
apiVersion: kubewise.io/v1alpha1
kind: RightSize
metadata:
  name: rightsize-conservative
  description: "Conservative right-sizing"
spec:
  percentile: p95          # p50, p90, p95, p99
  buffer: 20               # percentage headroom above the percentile
  scope:
    namespaces: ["*"]      # ["*"] = all, or specific list
    exclude:
      - namespace: kube-system
      - label: "kubewise.io/skip=true"
  limits:
    strategy: ratio        # ratio, fixed, unbounded
```

### Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `percentile` | string | p95 | Which usage percentile to base new requests on |
| `buffer` | int | 20 | Percentage headroom above the percentile |
| `scope.namespaces` | string[] | ["*"] | Namespaces to include |
| `scope.exclude` | object[] | [] | Namespace or label exclusions |
| `limits.strategy` | string | ratio | How to compute new limits |

### Limit strategies

- **ratio**: Maintains the original request-to-limit ratio. If the original request was 100m with a 200m limit (2x ratio), the new limit will be `newRequest * 2.0`.
- **fixed**: Sets the limit equal to the request (no burst allowed).
- **unbounded**: Removes limits entirely (set to 0).

### CLI equivalent

```bash
kubectl whatif rightsize --percentile=p95 --buffer=20 --exclude-namespace=kube-system --limit-strategy=ratio
```

## Consolidate

Replaces current node types with a target instance type and determines the minimum number of nodes needed.

### YAML schema

```yaml
apiVersion: kubewise.io/v1alpha1
kind: Consolidate
metadata:
  name: consolidate-m6i
  description: "Consolidate to m6i.xlarge"
spec:
  target_node_type: m6i.xlarge
  max_nodes: 50            # 0 = unlimited
  keep_node_pools:
    - gpu-pool
```

### Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `target_node_type` | string | required | Target instance type |
| `max_nodes` | int | 0 | Maximum number of nodes (0 = unlimited) |
| `keep_node_pools` | string[] | [] | Node pool names to leave untouched |

### How it works

1. Removes all nodes not in `keep_node_pools`
2. Creates virtual nodes of the target type
3. Runs bin-packing simulation to place all pods
4. Adds nodes until all pods are scheduled or `max_nodes` is hit
5. Reports the final node count and any unschedulable pods

### CLI equivalent

```bash
kubectl whatif consolidate --node-type=m6i.xlarge --max-nodes=50 --keep-pool=gpu-pool
```

## SpotMigrate

Identifies workloads eligible for spot/preemptible instances and estimates cost savings and eviction risk.

### YAML schema

```yaml
apiVersion: kubewise.io/v1alpha1
kind: SpotMigrate
metadata:
  name: spot-stateless
  description: "Move stateless workloads to spot"
spec:
  eligibility:
    min_replicas: 2
    controller_types:
      - Deployment
      - ReplicaSet
    exclude_namespaces:
      - kube-system
      - payments
  spot_discount: 0.65
```

### Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `eligibility.min_replicas` | int | 2 | Minimum replica count for eligibility |
| `eligibility.controller_types` | string[] | all | Controller types eligible for spot |
| `eligibility.exclude_namespaces` | string[] | [] | Namespaces to exclude |
| `spot_discount` | float | 0.65 | Spot discount fraction (0.65 = 65% off) |

### CLI equivalent

```bash
kubectl whatif spot --min-replicas=2 --discount=0.65 --exclude-namespace=kube-system,payments
```

## Composite

Chains multiple scenarios sequentially. Each step receives the output of the previous step.

### YAML schema

```yaml
apiVersion: kubewise.io/v1alpha1
kind: Composite
metadata:
  name: aggressive-savings
  description: "Right-size then move to spot"
spec:
  steps:
    - kind: RightSize
      spec:
        percentile: p95
        buffer: 20
    - kind: SpotMigrate
      spec:
        eligibility:
          min_replicas: 2
        spot_discount: 0.65
```

### Common use cases

**Conservative savings**: Right-size with wide buffers, keep everything on-demand.

```bash
kubectl whatif apply -f scenarios/rightsize-conservative.yaml
```

**Maximum savings**: Right-size aggressively, then move stateless workloads to spot.

```bash
kubectl whatif apply -f scenarios/composite-savings.yaml
```

**Compare approaches**: See the tradeoff between conservative and aggressive right-sizing.

```bash
kubectl whatif compare -f scenarios/rightsize-conservative.yaml -f scenarios/rightsize-aggressive.yaml
```
