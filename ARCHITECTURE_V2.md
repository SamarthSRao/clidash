# CLIDash V2: Jaeger-Inspired Sidecar Architecture

Following the **Jaeger/OpenTelemetry** pattern, CLIDash is evolving from a single CLI into a distributed orchestration system. This allows the RL Agent to manage consistency at scale without creating "dependency hell" for microservice developers.

## 🏗️ The 4-Tier Architecture

### 1. `clidash-sdk` (The Client)
- **Role**: Lightweight library embedded in the Microservice (Go, Python, Java).
- **Communication**: Sends telemetry (Request Type, Latency) to the local Agent via **UDP** (fire-and-forget, zero latency impact).
- **Control**: Intercepts DB calls to apply the `ConsistencyLevel` currently mandated by the Agent.

### 2. `clidash-agent` (The Sidecar)
- **Deployment**: Runs as a sidecar container (Kubernetes) or daemon on the same host as the microservice.
- **Role**: 
    - Deduplicates and batches metrics before sending to the Optimizer.
    - Maintains a local cache of the **Consistency Policy** (Strong/Eventual) received from the Optimizer.
    - Handles discovery and connection to the global Optimizer.

### 3. `clidash-optimizer` (The Collector & Brain)
- **Role**: The central "State Department."
    - Aggregates metrics from thousands of Agents.
    - Runs the **DQN RL Agent** to calculate optimal consistency levels per service/route.
    - Pushes policy updates back to Agents via **gRPC Stream** or Long Polling.

### 4. `clidash-dashboard` (The UI)
- **Role**: Connects to the Optimizer to provide global visibility.
- **View**: Real-time graph of system-wide throughput, latency saved, and AI decision logs.

---

## 🔄 Data Flow (The "Push" Model)

1. **Microservice** performs a "Payment" action. 
2. **SDK** sends a UDP packet to `localhost:5775`: `{ "svc": "payments", "op": "write", "latency": 45ms }`.
3. **Agent** batches this with 100 other requests and forwards it to the **Optimizer**.
4. **Optimizer** sees a latency spike on the DB and decides: *"Catalog service is causing lock contention. Switch Catalog to EVENTUAL."*
5. **Optimizer** pushes policy `{ "svc": "catalog", "mode": "EVENTUAL" }` to all **Catalog Agents**.
6. **Catalog Agents** update local cache; next Catalog DB read is instantly "Fast but Messy."

---

## 📁 New Directory Structure

```
/cmd
  /agent       # The Sidecar binary
  /optimizer   # The Central Brain binary
  /dashboard   # The TUI Dashboard binary
/pkg
  /proto       # Common message formats (Protobuf/JSON)
  /sdk         # Client libraries
/internal
  /rl          # The DQN logic
```
