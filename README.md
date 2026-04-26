# RateNote
RateNote is a Go project built using Clean Architecture principles.  
It provides a backend system for managing items with ratings, comments, and optional images.

The project uses PostgreSQL as a database and exposes HTTP handlers for working with data.  
A simple HTML interface is included, but the core focus is the backend architecture.

# Technologies
- Go 1.24+
- Docker & Docker Compose 3.9+

# Start Instructions
## 1 Clone Repo
```
git clone https://github.com/Xiancel/RateNote.git
cd RateNote
```
## 2 Run Docker
```
docker-compose --profile dev up --build

```

# Project diagram

Clean Architecture (simplified)
```
┌──────────────────────────────────────────────┐
│ Layer 4: Frameworks & Drivers               │
│----------------------------------------------│
│  - HTTP Server (net/http + chi)             │
│  - Docker / Env                             │
│  - PostgreSQL                               │
│  - main.go (DI)                             │
│                                              │
│  ┌──────────────────────────────────────────┐│
│  │ Layer 3: Interface Adapters              ││
│  │------------------------------------------││
│  │  HTTP Handlers                           ││
│  │   - ItemHandler (JSON API)               ││
│  │   - ItemPageHandler (HTML UI)            ││
│  │                                          ││
│  │  Templates (SSR UI)                      ││
│  │   - home.html                            ││
│  │   - item.html                            ││
│  │   - add.html                             ││
│  │   - edit.html                            ││
│  │                                          ││
│  │  Request/Response mapping                ││
│  │                                          ││
│  │  ┌──────────────────────────────────────┐││
│  │  │ Layer 2: Use Cases (Application)     │││
│  │  │--------------------------------------│││
│  │  │  ItemService                         │││
│  │  │   - AddItem                          │││
│  │  │   - GetItem                          │││
│  │  │   - ListItem                         │││
│  │  │   - UpdateItem                       │││
│  │  │   - DeleteItem                       │││
│  │  │                                      │││
│  │  │  Business logic                      │││
│  │  │                                      │││
│  │  │  ┌──────────────────────────────────┐│││
│  │  │  │ Layer 1: Domain (Core)          ││││
│  │  │  │----------------------------------││││
│  │  │  │  Entity: Item                    ││││
│  │  │  │   - id                           ││││
│  │  │  │   - name                         ││││
│  │  │  │   - rating                       ││││
│  │  │  │   - comment                      ││││
│  │  │  │   - image_path                   ││││
│  │  │  │                                  ││││
│  │  │  │  Interfaces                      ││││
│  │  │  │   - ItemRepository               ││││
│  │  │  └──────────────────────────────────┘│││
│  │  └──────────────────────────────────────┘││
│                                              ││
└──────────────────────────────────────────────┘│
│ External Resource                             │
│   - PostgreSQL                                │
└──────────────────────────────────────────────┘
```
