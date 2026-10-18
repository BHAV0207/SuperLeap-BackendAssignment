# Runbook

What "healthy" means, what each alert means and how to check and fix it.

## Healthy looks like

- `kubectl -n lead-crm get pods`: 2/2 api pods `Running` and ready, `postgres-0` and `redis` ready
- `curl localhost:8088/readyz` → `{"checks":{"postgres":"up","redis":"up"},"success":true}`
- Grafana **Lead CRM** dashboard: availability ≥ 99.5%, p95 < 300ms, both dependencies at 1
- No `LeadCrm*` alerts at http://localhost:9090/alerts
- `make smoke` passes

## SLOs

| SLI | Target | Measured by |
|---|---|---|
| Availability | 99.5% of user requests are not 5xx | `leadcrm:error_ratio:rate5m` |
| Latency | 95% of user requests < 300ms | `leadcrm:latency_p95:rate5m` |

Health/readiness/metrics endpoints are excluded from both.

## Alerts

### LeadCrmHighErrorRate

More than 5% of requests return 5xx for 5 minutes (critical).

1. Dashboard: which route has 5xx? Did it start right after a deploy?
   (`helm -n lead-crm history lead-crm`)
2. Logs for failing requests:
   `kubectl -n lead-crm logs -l app.kubernetes.io/name=lead-crm --prefix | grep '"level":"ERROR"'`
3. Take a `trace_id` from a log line into Jaeger to see which call failed.
4. If it started with a deploy: `make rollback`, then confirm the error ratio
   drops back to 0 on the dashboard and `make smoke` passes.
5. If not deploy related, check `LeadCrmPostgresDown` / db pod.

### LeadCrmHighLatency

p95 above 300ms for 10 minutes (warning).

1. Dashboard latency panel: all routes or one?
2. Cache panel: lots of `error`/`miss`? Redis may be down (reads all go to Postgres).
3. Open a slow trace in Jaeger, compare the SQL span with the total.
4. Check throttling / memory: `kubectl -n lead-crm top pods` (needs metrics-server) or the pod CPU graphs in the default Grafana Kubernetes dashboards.

### LeadCrmNoReadyPods

No api pod is ready for 2 minutes (critical, the API is down).

1. `kubectl -n lead-crm get pods` and `kubectl -n lead-crm describe pod <pod>`
2. `ImagePullBackOff` / `ErrImageNeverPull` → wrong tag, `make rollback`
3. `CrashLoopBackOff` → `kubectl -n lead-crm logs <pod> --previous`
4. Running but not ready → `curl localhost:8088/readyz` usually shows Postgres down

### LeadCrmPostgresDown

API pods cannot reach PostgreSQL for 1 minute (critical).

1. `kubectl -n lead-crm get pod postgres-0`, `kubectl -n lead-crm logs postgres-0`
2. PVC bound? `kubectl -n lead-crm get pvc`
3. Was a NetworkPolicy changed? `kubectl -n lead-crm get networkpolicy`
4. Once postgres is ready again the api pods become ready on their own, no restart needed.

### LeadCrmCacheDown

Redis unreachable for 5 minutes (warning). Users are not affected apart from
slightly slower `GET /leads/:id`.

1. `kubectl -n lead-crm get pods -l app.kubernetes.io/name=redis`
2. The api reconnects automatically when Redis is back; cache `hit` count on the dashboard starts rising again.

## Fault demo: Redis outage

Shows the cache failure boundary and recovery. Run with the stack up.

```bash
# baseline
curl -s localhost:8088/readyz
# break it
kubectl -n lead-crm scale deploy/redis --replicas=0
# API still works, readiness reports redis degraded
curl -s -w ' %{http_code} %{time_total}s\n' localhost:8088/leads
curl -s localhost:8088/readyz
# recover
kubectl -n lead-crm scale deploy/redis --replicas=1
```

- **Symptom:** `dependency_up{dependency="redis"}` drops to 0, cache lookups show `error`.
- **Hypothesis:** Redis is gone, the api falls back to Postgres.
- **Evidence:** `/readyz` says `redis: degraded`, requests still 200 in a few ms, `LeadCrmCacheDown` goes pending then firing after 5m.
- **Fix:** scale Redis back.
- **Verify:** `/readyz` shows `redis: up`, cache `hit` rate climbs again, alert resolves.

## Fault demo: bad release

```bash
# deploy an image tag that does not exist
helm upgrade lead-crm deploy/helm/lead-crm -n lead-crm \
  --reuse-values --set app.image.tag=badbad0 --atomic --wait --timeout 60s
helm -n lead-crm history lead-crm --max 3
```

- **Symptom:** new pod stuck in `ImagePullBackOff`, rollout does not progress.
- **Evidence:** old pods are never removed (`maxUnavailable: 0`), so traffic keeps flowing.
- **Fix:** `--atomic` rolls back automatically after the timeout; history shows `Rollback to <n>`.
- **Verify:** `curl localhost:8088/version` shows the previous tag, `make smoke` passes.
