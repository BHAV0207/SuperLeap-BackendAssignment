# Design decisions

Short notes on why things are set up the way they are, what else I
considered, and when I would choose differently.

## Environment: local kind cluster

- **Why:** free, reproducible with one config file, close enough to a real
  cluster (real scheduler, Services, NetworkPolicy, PVCs).
- **Alternatives:** Minikube (similar), a managed cluster like EKS/GKE (costs
  money, needs Terraform or click-ops).
- **Would change if:** other people needed to reach it, or I needed node
  failure tolerance. Then a managed cluster with an Ingress and LoadBalancer.

## Workloads

| Component | Choice | Why | Alternative |
|---|---|---|---|
| API | Deployment, 2 replicas | Stateless, any pod can serve any request. 2 so one can restart or roll without downtime | 1 replica (downtime on every deploy), HPA (no load to justify it) |
| PostgreSQL | StatefulSet, 1 replica, PVC | Needs stable identity and storage that survives restarts | Managed db (best for real data), operator like CloudNativePG (replication + backups, too much for this) |
| Redis | Deployment, no PVC | Only a read cache with a 10 min TTL, losing it loses nothing | StatefulSet with AOF (only if Redis held real data) |

## Health checks

- **Liveness `/healthz`** only checks the process. If it checked Postgres, a
  db outage would make Kubernetes restart every api pod over and over, which
  makes things worse and does not fix the db.
- **Readiness `/readyz`** checks Postgres (required). Redis is reported but
  does not fail readiness because the app works without it.
- **Startup probe** gives up to 60s on boot because the app retries the db
  connection for ~30s. Without it, liveness could kill a pod that is just
  waiting for Postgres.
- **preStop sleep 5s** + graceful shutdown on SIGTERM: the pod keeps serving
  while it is removed from the Service endpoints, then finishes in-flight
  requests. This is why a rolling update shows zero failed requests.

## Resources

Requests are set close to what the pods actually use at idle so scheduling is
honest. Memory limits stop a leak from taking down the node. The api CPU limit
is generous (500m) because Go handles CPU throttling badly when limits are
tight.

## Configuration: one Helm chart

- **Why Helm:** one templated chart instead of copies of YAML per
  environment. Values hold everything that changes (tag, replicas,
  resources). The template refuses to render without an image tag, so
  `latest` can't sneak in.
- **Alternative:** Kustomize overlays. Also fine; Helm gives release history
  and `rollback`, which I use for release control.
- **Secrets are not in the chart.** The chart references an existing Secret.

## CI/CD

- **GitHub Actions** because the code is already on GitHub and it needs no
  extra servers.
- **Image tag = commit sha.** Tags are immutable in practice and map 1:1 to a
  commit. The same sha is baked into the binary (`/version`, logs, metrics).
- **Build once, deploy the same image.** CD pulls the image CI scanned; it
  never rebuilds.
- **Why a self-hosted runner instead of Argo CD:** the cluster is on my
  laptop, so something on the laptop has to apply changes. A self-hosted
  runner is one small process; Argo CD would mean another controller, a
  repo for desired state and a way to bump image tags in git. On a shared
  real cluster I would pick Argo CD (pull based, drift detection, the cluster
  doesn't need to trust CI).
- **Manual approval on deploy** through a GitHub environment: CI is
  automatic, putting a build on the running cluster is a human decision. Also
  stops random workflow runs from executing on my machine.

## Release strategy: rolling update + automatic rollback

- **Why:** the app is stateless and the DB schema change is additive
  (GORM AutoMigrate only adds), so old and new versions can run side by side
  for a few seconds. Rolling with `maxUnavailable: 0` gives zero downtime with
  no extra tooling.
- `helm --atomic` rolls back automatically when new pods never get ready.
  `make rollback` covers the case where pods are ready but the behaviour is
  wrong (caught by the smoke test or the error rate alert).
- **Alternatives:** blue-green (2x resources, instant switch back; worth it
  if a rollout needed to be atomic across many pods), canary with
  Argo Rollouts / a service mesh (send 10% of traffic first; worth it with
  real user traffic to compare error rates). Feature flags for risky changes
  in code.
- **Would change if:** a release contained a breaking DB migration. Then I
  would split it into expand → deploy → contract steps, since rolling back the
  app would not roll back the schema.

## Observability

- **Prometheus + Grafana (kube-prometheus-stack):** standard, gives
  kube-state-metrics (pod restarts, ready replicas) for free, and alerts as
  code with PrometheusRule.
- **Route label uses the template (`/leads/:id`)** not the raw path,
  otherwise every lead id becomes a new time series.
- **Probe traffic is excluded from SLIs**, it would make availability look
  better than what users get.
- **OpenTelemetry → Jaeger:** vendor neutral SDK, the same code could send to
  Tempo or a SaaS. The trace showed directly where time goes (Redis vs SQL).
- **No Loki:** JSON logs + `kubectl logs` + trace_id was enough for one
  service. With several services or long retention I would add Loki.

## Security

- Distroless non-root image (no shell to exploit, 60MB, 0 HIGH/CRITICAL CVEs
  at build).
- Default-deny NetworkPolicy so a compromised pod elsewhere can't reach the db.
- No ServiceAccount token: the app never calls the Kubernetes API.
- Actions pinned to commit shas: a moved tag can't inject code into the
  pipeline (this happened to a popular security action in 2026).
- govulncheck caught a real issue: the original project's Postgres driver
  (pgx 5.6.0) had a known vulnerability. Upgraded to 5.11.0.

## Bugs found while doing this

1. Redis outage added **3.7s** to every request: go-redis retried the dial 5
   times with default timeouts. Set short timeouts and 1 dial attempt → ~8ms.
2. If Redis was down at startup, caching stayed off until the app restarted
   (the client was set to nil). Now go-redis reconnects by itself.
3. `GET /leads?status=CONVERTED` returned 400, the status validator was
   missing CONVERTED.
4. `.gitignore` ignored `go.mod` and `go.sum`, so builds were not
   reproducible.
