# Arcade CATE + RCF Integration Mock

This harness demonstrates how Arcade's CATE (Contextual Access) and RCF (Runtime Consequence Firewall) operate in a unified post-auth pipeline. CATE handles identity brokering, and RCF hooks the execution phase to pre-detonate payloads in a `<1ms` wazero micro-VM.

## Architecture Flow

```mermaid
sequenceDiagram
    participant LLM as Agent/LLM
    participant CATE as Arcade CATE (Auth)
    participant RCF as RCF (Micro-VM)
    participant API as Target System

    LLM->>CATE: Request Tool Execution
    CATE-->>CATE: Validate OAuth / Identity
    CATE->>RCF: Handoff Payload + user_id
    RCF-->>RCF: Prune Schema (0.3ms BM25)
    RCF-->>RCF: Detonate in wazero micro-VM
    alt Malicious (e.g. DROP TABLE)
        RCF->>CATE: JSON-RPC Error (Denied)
        RCF-->>RCF: Write WORM Log + user_id
    else Safe
        RCF->>API: Execute Tool
    end
```

## Expected Output
*(Simulating Safe Execution, Malicious Injection, and Missing Identity Context)*

![CATE + RCF Terminal Output](./terminal_output.png)
