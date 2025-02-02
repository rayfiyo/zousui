# zousui

- 文明進化シミュレーター
- Civilization Evolution Simulator

## .env

`backend/.env`

```
GEMINI_API_KEY=
OPENAI_API_KEY=
```

## tree

### backend

```bash
tree backend/
```

```text
backend/
├── cmd
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
│   │   ├── image_controller.go
│   │   └── simulate_controller.go
│   ├── gateway
│   │   ├── gemini_gateway.go
│   │   └── moc_gateway.go
│   └── router
│       └── router.go
├── usecase
│   ├── community.go
│   ├── diplomacy.go
│   └── simulate.go
└── utils
    ├── config
    │   └── env.go
    └── consts
        └── consts.go
```

### frontend

```bash
tree -I "node_modules" frontend/
```

```text
frontend/
├── app
│   ├── community
│   │   ├── [id]
│   │   │   └── page.tsx
│   │   └── new
│   │       └── page.tsx
│   ├── diplomacy
│   │   └── page.tsx
│   ├── layout.tsx
│   ├── NavBar.tsx
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
