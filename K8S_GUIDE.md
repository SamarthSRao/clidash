# Kubernetes Deployment Guide for CLIDash

Deploying CLIDash in Kubernetes (K8s) leverages the **Sidecar Container Pattern**, exactly like Jaeger and Istio.

## 1. Architecture in K8s

- **Central Optimizer**: Deployed as a single-replica `Deployment` with a `ClusterIP` Service. This is the global "Brain."
- **Sidecar Agent**: Injected into every Application Pod. It shares the same network namespace (`localhost`) as your app.

## 2. Steps to Deploy

### Step A: Build Images
```bash
docker build -t your-reg/clidash-optimizer:latest --target optimizer .
docker build -t your-reg/clidash-agent:latest --target agent .
```

### Step B: Deploy the Brain
Apply the optimizer manifest:
```bash
kubectl apply -f k8s/optimizer.yaml
```
The brain will be reachable inside the cluster at `http://clidash-optimizer.default.svc.cluster.local`.

### Step C: Inject the Sidecar
Modify your application's Deployment to include the `clidash-sidecar` container. See `k8s/app-sidecar-example.yaml` for a full template.

Key points:
1.  **Shared Network**: Your app talks to the agent via `localhost:5775`.
2.  **Environment Variables**:
    *   `OPTIMIZER_URL`: Points to the internal K8s service.
    *   `SERVICE_ID`: Can be mapped to the Pod Name using the K8s Downward API.

## 3. Benefits of this Setup
- **Zero Configuration**: The app always knows the agent is at `localhost`.
- **Fault Tolerance**: If one agent fails, only that pod is affected. The global site stays up.
- **Dynamic Scaling**: As you scale your app replicas (`kubectl scale --replicas=10`), the RL Optimizer automatically sees 10 new agents and includes them in the rewards calculation.
