VanosDNS

VanosDNS is a distributed, mesh-based Domain Name System designed to maximize resilience, performance, and operational independence.
It aims to eliminate central points of failure while providing enterprise-grade DNS resolution with extreme tail-latency optimization.

VanosDNS is built as a community-driven infrastructure: free to use, transparent by design, and protected against commercial exploitation.

Objectives
    VanosDNS addresses structural limitations of traditional DNS infrastructures:
    Elimination of central control points
    Resistance to large-scale outages and DDoS attacks
    Reduction of geopolitical and provider lock-in risks
    Optimization of real-world latency, including worst-case scenarios
    Operation without monetization, data exploitation, or traffic resale

Key Features (Version 1)
    Distributed Mesh Network
    Each user instance can operate as a network node
    Local and regional peer discovery
    No centralized authority or root infrastructure
    Intelligent DNS Resolution
    Dynamic node selection based on latency and reliability
    Parallel resolution and adaptive fallback strategies
    Request cancellation to reduce tail latency impact
    Advanced Anti-Sybil Mechanisms
    No proof-of-work, tokens, or economic incentives
    Behavioral and network-based node scoring
    Progressive isolation of unreliable or artificial nodes
    Caching Architecture
    Local cache for immediate performance gains
    Distributed cache sharing between trusted peers
    Automatic cache invalidation and consistency mechanisms
    Observability and Metrics
    Local metrics collection (latency, error rate, cache efficiency)
    No centralized logging
    No data collection, resale, or traffic analysis
    Local Administration Interface
    Local GUI for node configuration and monitoring
    No cloud dependency
    No account, authentication service, or external control plane

Architecture Overview : 
  DNS Client
     |
     v
  Local Resolver
     |
     +-- Cache Layer
     |
     +-- Routing Engine
           |
           +-- Local High-Score Nodes
           +-- Geographically Redundant Peers


Routing decisions are continuously adapted based on real-time network conditions and historical node behavior.

Repository Structure : 

vanosdns/
├── cmd/
│   └── vanosdns/
│       └── main.go
├── internal/
│   ├── resolver/
│   │   ├── resolver.go
│   │   ├── cache.go
│   │   └── routing.go
│   ├── config/
│   │   └── config.go
│   ├── metrics/
│   │   └── metrics.go
│   └── types/
│       └── types.go
├── gui/
│   └── (local administration interface)
├── docs/
│   └── architecture.md
├── LICENSE
└── README.md

Installation
Version 1 is production-capable but actively evolving.
    git clone https://github.com/vanosdns/vanosdns.git
    cd vanosdns
    go build ./cmd/vanosdns
    ./vanosdns

Enterprise Usage :

Organizations and enterprises are explicitly allowed to:
    Deploy VanosDNS in production environments
    Use it for internal or external DNS resolution
    Integrate it into existing infrastructure stacks

Restrictions:
    VanosDNS may not be sold, licensed, or monetized
    It may not be included in paid products or services
    DNS traffic and metadata may not be monetized
    See the LICENSE file for full terms.

License
    VanosDNS is distributed under the VanosDNS Community License 1.2.
    Free to use
    Free to modify
    Free to deploy (including enterprises)
    Strict prohibition of monetization or commercial resale

Contributions : 
    Contributions are welcome and encouraged.
    This includes:
    Performance improvements
    Security analysis
    Networking optimizations
    Documentation and architecture reviews

All contributions must comply with the project license.

Roadmap :
    V1 – Distributed resolver, mesh routing, anti-Sybil, local GUI
    V1.5 – Performance tuning based on community feedback
    V2 – Advanced routing strategies and network intelligence
    V2.5 – Community-driven improvements
    V3 – Long-term scalability and protocol evolution

Philosophy and VanosDNS prioritizes:
    Resilience over convenience
    Transparency over monetization

Infrastructure integrity over market capture

It is designed to be a foundational component of a more decentralized and reliable internet.
