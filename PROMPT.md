# AI Implementation Prompt

Use the following prompt with an LLM (Claude/GPT) to further develop the AI logic for this system:

---

**Prompt for RL Agent Development:**

"I am building a Reinforcement Learning-Based Dynamic Consistency Optimizer for an e-commerce microservices architecture. 

The goal is to implement a Deep Q-Network (DQN) agent that acts as a 'Traffic Controller.' It needs to observe:
1. Current throughput (Requests per second)
2. p99 Latency per service
3. Business criticality of the operation (High for Checkout/Payment, Low for Views)
4. Current resource saturation (CPU/Memory/Lock wait times)

The Agent has three actions per service:
- `SET_STRONG`: Prioritize data integrity (Slow but Perfect).
- `SET_EVENTUAL`: Prioritize performance (Fast but Messy).
- `SET_SESSION`: Hybrid approach.

The reward function should be:
`Reward = (Latency_Target - Current_Latency) * Weight_Perf - (Data_Correction_Cost) * Weight_Integrity`

Please provide a Python implementation using PyTorch and Gymnasium for this environment, including a mock environment that simulates e-commerce 'Black Friday' traffic spikes and the corresponding 'lock-wait' cascades that happen when many Strong consistency reads hit a hot-key."

---
