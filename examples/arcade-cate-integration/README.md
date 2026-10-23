{\rtf1\ansi\ansicpg1252\cocoartf2907
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fswiss\fcharset0 Helvetica;}
{\colortbl;\red255\green255\blue255;}
{\*\expandedcolortbl;;}
\paperw11900\paperh16840\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\pard\tx720\tx1440\tx2160\tx2880\tx3600\tx4320\tx5040\tx5760\tx6480\tx7200\tx7920\tx8640\pardirnatural\partightenfactor0

\f0\fs24 \cf0 # Arcade CATE + RCF Integration Mock\
\
This harness demonstrates how Arcade's CATE (Contextual Access) and RCF (Runtime Consequence Firewall) operate in a unified post-auth pipeline. CATE handles identity brokering, and RCF hooks the execution phase to pre-detonate payloads in a `<1ms` wazero micro-VM.\
\
## Architecture Flow\
\
```mermaid\
sequenceDiagram\
    participant LLM as Agent/LLM\
    participant CATE as Arcade CATE (Auth)\
    participant RCF as RCF (Micro-VM)\
    participant API as Target System\
\
    LLM->>CATE: Request Tool Execution\
    CATE-->>CATE: Validate OAuth / Identity\
    CATE->>RCF: Handoff Payload + user_id\
    RCF-->>RCF: Prune Schema (0.3ms BM25)\
    RCF-->>RCF: Detonate in wazero micro-VM\
    alt Malicious (e.g. DROP TABLE)\
        RCF->>CATE: JSON-RPC Error (Denied)\
        RCF-->>RCF: Write WORM Log + user_id\
    else Safe\
        RCF->>API: Execute Tool\
    end\
}