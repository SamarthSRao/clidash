# CLIDash Architecture & Roadmap

## 1. System Structure

The `clidash` tool is designed as a **Centralized Orchestrator** for consistency management across a distributed environment.

### Core Modules:
- **Metrics Aggregator**: Collects real-time traffic (RPS), latency (p99), and error rates from registered microservices via a sidecar pattern or middleware.
- **RL Agent (DQN Brain)**: A reinforcement learning engine that calculates the "Reward" (Latency Improvement vs. Data Integrity Risk). It predicts if a "Strong" consistency request will cause a cascade failure and opts for "Eventual" instead for non-critical paths.
- **Protocol Driver**: Interacts with specific database APIs:
  - **DynamoDB Driver**: Toggles `ConsistentRead` flag.
  - **Cosmos DB Driver**: Switches consistency level via session tokens or header overrides.
  - **Cassandra Driver**: Adjusts `ConsistencyLevel` (ONE/QUORUM/ALL).

## 2. Dynamic Consistency Logic
| Load Level | Transaction Type | AI Decision | Business Impact |
| --- | --- | --- | --- |
| Low | Any | Strong | 100% Accuracy, negligible latency impact. |
| High | Payment / Inventory | Strong | Integrity priority. The AI throttles non-criticals to protect this. |
| High | Product View / Analytics | Eventual | 85% Latency reduction. Prevents site-wide lag. |
| Spiking | Cart / Session | Session | Balanced approach for UX continuity. |

## 3. Frontend Web Integration (Next Step)
While the CLI provides low-level control and real-time visualization, a web dashboard would offer:
- **Business Rule Editor**: Drag-and-drop criticality setting.
- **Cost/Savings Calculator**: Shows $ saved by hardware optimization and reduced "Oversell" errors.
- **Architectural Advisor**: Recommends consistency models based on historical traffic patterns.

## 4. How to Run `clidash`
Currently, the tool runs in **Simulation Mode** to demonstrate the RL behavior.

```bash
./clidash.exe
```
(Press 'q' to quit)
