# VanosDNS - Decentralized P2P DNS

VanosDNS is a decentralized DNS resolver based on a distributed trust architecture.
Each node acts as a client and potentially as a "bankroot" (coordination node).

## Features
- **Zero-Config:** TCP hole punching to traverse NAT without opening ports.
- **Meritocracy:** Dynamic peer scoring (Latency * Success / Age).
- **Security:** Identity based on Ed25519 and majority consensus (2/3).
- **Hybrid:** RAM cache (L1) / SQLite SSD (L2) and fallback root hints.
