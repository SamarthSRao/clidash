# Integration Guide: Connecting Services to CLIDash

To connect your microservices to the CLIDash RL Optimizer, follow the **Jaeger Sidecar Pattern**.

## 1. Deploy the Sidecar
Every microservice instance should have a `clidash-agent` running alongside it (on the same host or in the same Kubernetes pod).

**Why?** This ensures the microservice only talks to `localhost`, keeping latency near zero and preventing "dependency hell."

```bash
# Set the unique ID for this service
$env:SERVICE_ID="my-service-alpha"
go run cmd/agent/main.go
```

## 2. Use the SDK (Go Example)
Import the `clidash/pkg/sdk` into your service to start reporting metrics and receiving AI-driven consistency policies.

```go
import "clidash/pkg/sdk"

func main() {
    // 1. Initialize the SDK (points to the local agent sidecar)
    client := sdk.NewClient("my-service-alpha", "localhost:5775")

    // 2. Wrap your database calls
    start := time.Now()
    
    // Check what consistency the AI mandates right now
    mode := client.GetConsistency() // returns "STRONG" or "EVENTUAL"
    
    db.Query("SELECT * FROM inventory", mode)
    
    // 3. Report the result back to the AI
    latency := time.Since(start).Milliseconds()
    client.RecordOperation("GET_INVENTORY", float64(latency))
}
```

## 3. Connecting Other Languages (Python/Node.js)
Since the Sidecar Agent listens on standard HTTP/UDP, you don't even need a Go SDK. You can connect any service using simple JSON over HTTP:

### **To Report Metrics (Push):**
**POST** `http://localhost:5775/metrics`
```json
{
  "service_id": "python-api",
  "operation": "user_login",
  "latency_ms": 12.5,
  "rps": 450
}
```

### **To Get Consistency Policy (Pull):**
**GET** `http://localhost:5775/policy`
```json
{
  "service_id": "python-api",
  "consistency": "EVENTUAL"
}
```

---

## 4. Global Discovery
The `clidash-agent` needs to know where the central **Optimizer** is. 
You can configure this via Environment Variables when starting the agent:
- `OPTIMIZER_URL`: The address of the central brain (default: `localhost:8080`).
