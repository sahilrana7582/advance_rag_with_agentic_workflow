# Sentinel AI — Autonomous SRE Engineer

> Automate incident response and cut MTTR by 70%. An agentic AI SRE that investigates, diagnoses, remediates, and verifies Kubernetes incidents — with human-in-the-loop approval at every destructive step.

## Architecture

| Layer | Technology | Responsibility |
|-------|-----------|----------------|
| Infrastructure tools | Go + client-go | Query Kubernetes, Prometheus, Loki |
| Agent orchestration | Python + LangGraph | Incident lifecycle state machine |
| LLM | Gemini (provider-abstraction) | Reasoning & root cause |
| Metrics | Prometheus | Application & infra metrics |
| Logs | Loki | Application log aggregation |
| Database | PostgreSQL | Incident state & history |
| Cache/Coordination | Redis | Job state & event coordination |
| Local K8s | Kind | Local Kubernetes cluster |
| API | REST (Go) | Alert ingestion & incident management |

## Folder Structure

```
sentinel-ai/
├── cmd/                    # Go entry points (sre-api, sre-worker, k8s-agent)
├── internal/               # Go packages (incident, tools, remediation, ...)
├── agent-orchestrator/     # Python agents + LangGraph orchestration
├── deployments/            # K8s manifests (kind, app, monitoring, sentinel)
├── test-infra/             # Sample app, failure scenarios, load tests
├── api/                    # OpenAPI spec
└── docs/                   # Architecture, incident lifecycle, security
```

## Quick Start

```bash
# 0. Prerequisites: go, python3, docker, kubectl, kind

# 1. Start local infrastructure (Postgres, Redis)
make infra-up

# 2. Create Kind cluster + deploy sample app, Prometheus, Loki
make kind-up

# 3. Build & run the Go API
make run-api

# 4. Run the Python orchestrator
cd agent-orchestrator
make run

# 5. Trigger an incident and watch Sentinel respond
make trigger-incident
```

## Milestones

1. **M0** — Repo scaffolding & local infra
2. **M1** — Kind cluster, sample app, Prometheus, Loki
3. **M2** — Go infrastructure tools (health, logs, metrics, deployments, k8s)
4. **M3** — Python orchestrator + 6 agents (LangGraph)
5. **M4** — Remediation + human approval + verification
6. **M5** — Dockerize, deploy to Kind, end-to-end demo

See `docs/architecture.md` for the full design.

## Safety Rules (Production-Grade from Day 1)

1. **Read-only by default** — investigation agents cannot modify infrastructure.
2. **Human approval** — every destructive remediation requires explicit approval.
3. **RBAC** — least-privilege Kubernetes service account.
4. **Audit trail** — every tool call, args, result, and identity recorded.
5. **Idempotency** — prevent duplicate remediation from repeated alerts.
6. **Evidence-based RCA** — root cause references actual logs, metrics, events, deployments.
7. **Verification** — recovery confirmed via actual application health, not just `kubectl`.