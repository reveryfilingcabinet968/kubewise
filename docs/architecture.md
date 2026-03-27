# Architecture

## Data flow

```
                         ┌──────────────┐
                         │  Kubernetes  │
                         │   API Server │
                         └──────┬───────┘
                                │
                         ┌──────▼───────┐
                         │  Collector   │  pkg/collector
                         │  (snapshot)  │
                         └──────┬───────┘
                                │
                    ┌───────────▼───────────┐
                    │   ClusterSnapshot     │  pkg/collector/types
                    │   (in-memory struct)  │
                    └───────────┬───────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Scenario Engine     │  pkg/scenario
                    │   (deep copy + mutate)│
                    └───────────┬───────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Simulator           │  pkg/simulator
                    │   (bin-pack + cost)   │
                    └───────────┬───────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Risk Scorer         │  pkg/risk
                    │   (OOM + eviction)    │
                    └───────────┬───────────┘
                                │
                    ┌───────────▼───────────┐
                    │   Output Renderer     │  pkg/output
                    │   (table/JSON/MD)     │
                    └───────────────────────┘
```

## Module responsibilities

### pkg/collector

Reads the full schedulable state of the cluster. Produces a typed `ClusterSnapshot` struct that every other module operates on.

- **snapshot.go**: Orchestrates collection from Kubernetes API (nodes, pods, controllers, HPAs, PDBs, PVCs)
- **prometheus.go**: Queries Prometheus for historical usage percentiles (P50/P90/P95/P99)
- **types.go**: All snapshot struct definitions (the data contract)
- **deepcopy.go**: Explicit deep copy methods for all types

### pkg/pricing

Maps instance types to hourly costs. Supports AWS, GCP, and Azure.

- Auto-detects cloud provider from node labels
- Caches pricing in `~/.kubewise/pricing/` with 24-hour TTL
- Falls back to user-provided YAML pricing file
- GCP prices CPU and memory separately; AWS and Azure price by instance type

### pkg/scenario

Pure mutations applied to snapshot copies. Three scenario types plus composites.

- **engine.go**: Deep copies the snapshot, applies the scenario, returns the mutated copy
- **rightsize.go**: Adjusts requests/limits based on usage percentiles
- **consolidate.go**: Replaces nodes with a target type
- **spot.go**: Tags eligible pods for spot scheduling
- **parser.go**: Parses scenario YAML files

### pkg/simulator

Scheduling simulation and cost calculation.

- **predicates.go**: Filter functions (FitsResources, MatchesAffinity, ToleratesTaints, SatisfiesTopologySpread)
- **priorities.go**: MostRequestedPriority scoring (pack tight, not spread)
- **binpack.go**: Core bin-packing loop (DaemonSet pre-placement, first-fit-decreasing)
- **autoscaler.go**: Adds virtual nodes until all pods fit
- **cost.go**: Monthly cost calculation with per-namespace allocation

### pkg/risk

Risk scoring for simulation outcomes.

- **oom.go**: OOM probability from usage percentiles vs new memory limits
- **eviction.go**: Spot interruption risk (`interruptionRate ^ replicaCount`)
- **scheduling.go**: Fraction of unschedulable pods
- **aggregate.go**: Cluster-wide rollup with green/amber/red classification

### pkg/output

Renders reports to terminal, JSON, or Markdown.

- **table.go**: Rich terminal output with lipgloss styling
- **json.go**: Stable, indented JSON with sorted keys
- **markdown.go**: GitHub PR comment format with collapsible sections

## Bin-packer algorithm

The simulator implements a simplified Kubernetes scheduler:

1. **Pre-place DaemonSet pods** on every matching node (respects taints and affinity)
2. **Sort remaining pods** by total resource request descending (first-fit-decreasing)
3. For each pod:
   - **Filter**: check FitsResources, MatchesAffinity, ToleratesTaints, SatisfiesTopologySpread
   - **Score**: MostRequestedPriority (prefer already-packed nodes)
   - **Place** on highest-scoring node
4. Unschedulable pods are collected for reporting

What we implement (MVP):
- Resource fit (CPU + memory)
- Required node affinity
- Taint tolerations
- Topology spread constraints (DoNotSchedule only)
- MostRequestedPriority + balance penalty scoring

What we skip (Phase 2):
- Pod preemption and priority classes
- Inter-pod anti-affinity
- Volume topology constraints
- Preferred scheduling terms

## Risk scoring

### OOM risk

Estimated from usage percentiles vs new memory limits using linear interpolation:

| Condition             | Estimated risk |
|-----------------------|----------------|
| newLimit > P99        | < 1% (green)   |
| P95 < newLimit <= P99 | 1-5% (amber)   |
| P90 < newLimit <= P95 | 5-10% (red)    |
| P50 < newLimit <= P90 | 10-50% (red)   |
| newLimit <= P50       | 50%+ (red)     |

Per-workload risk: `1 - product(1 - container_risk)` across all containers.

### Spot eviction risk

Based on historical interruption rates by instance family:

- Risk of ALL replicas interrupted simultaneously: `interruptionRate ^ replicaCount`
- Green: < 0.1%, Amber: 0.1-1%, Red: > 1%

### Scheduling risk

- Fraction of pods that cannot be placed: `unschedulable / total`
- Green: 0%, Amber: < 1%, Red: >= 1%

## Pricing data

| Provider | Source                 | Pricing model         |
|----------|------------------------|-----------------------|
| AWS      | EC2 Pricing API (Bulk) | Per instance type     |
| GCP      | Cloud Billing Catalog  | Per vCPU + per GB RAM |
| Azure    | Retail Prices API      | Per instance type     |

Pricing is cached locally in `~/.kubewise/pricing/{provider}_{region}.json` with a 24-hour TTL. If cloud APIs are unreachable, users can provide a YAML pricing file with `--pricing-file`.
