# Silvercord on Kubernetes

Manifests for deploying Silvercord (Discord bot + Go API + Python agent) onto a k3s cluster.

## Layout

- `namespace.yaml` — `silvercord` namespace
- `agent.yaml` — Flask agent + ChromaDB PVC + ClusterIP service (Recreate strategy due to RWO volume)
- `api.yaml` — Go API + ClusterIP service (RollingUpdate w/ `/ping` probes, zero-downtime)
- `bot.yaml` — Discord bot (Recreate; gateway forbids duplicate sessions)
- `kustomization.yaml` — applies all of the above

## Secrets

Three secrets are required, loaded from environment files (kept out of git):

```sh
kubectl -n silvercord create secret generic bot-env --from-env-file=../.env
kubectl -n silvercord create secret generic api-env --from-env-file=../api/.env.local
kubectl -n silvercord create secret generic agent-env --from-env-file=../silvercord_agent/.env
```

The bot's `API_BASE_URL` and `AGENT_API_URL` are set in `bot.yaml` directly (cluster-internal DNS), so omit those from `.env` or they'll be overridden by the explicit `env:` entries anyway.

## Apply

```sh
kubectl apply -k .
```

Image tags are bumped by the deploy workflow via `kubectl set image`.
