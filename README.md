# zousui

- 文明進化シミュレーター
- Civilization Evolution Simulator

## tree

### backend

```bash
tree backend/
```

```bash
backend/
├── cmd
│   ├── cmd
│   └── main.go
├── domain
│   ├── entity
│   │   ├── agent.go
│   │   └── community.go
│   └── repository
│       └── simulate.go
├── go.mod
├── go.sum
├── infrastructure
│   └── repository
│       ├── memory_agent_repo.go
│       └── memory_community_repo.go
├── interface
│   ├── controller
│   │   ├── community_controller.go
│   │   ├── diplomacy_controller.go
│   │   └── simulate_controller.go
│   ├── gateway
│   │   └── llm_gateway.go
│   └── router
│       └── router.go
├── usecase
│   ├── community.go
│   ├── diplomacy.go
│   └── simulate.go
└── utils
    └── const.go
```

### frontend

```bash
tree -I "node_modules" frontend/
```

```bash
frontend/
├── app
│   ├── layout.tsx
│   └── page.tsx
├── eslint.config.mjs
├── next.config.ts
├── next-env.d.ts
├── package.json
├── package-lock.json
├── public
│   ├── file.svg
│   ├── globe.svg
│   ├── next.svg
│   ├── vercel.svg
│   └── window.svg
├── README.md
├── src
│   └── app
│       ├── favicon.ico
│       ├── globals.css
│       ├── layout.tsx
│       ├── page.module.css
│       └── page.tsx
└── tsconfig.json
```
