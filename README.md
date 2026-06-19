# VK Jeopardy

A **Jeopardy-style quiz game for [VKontakte](https://vk.com) community chats**,
with an admin/rating web app built as a VK Mini App.

This is a monorepo:

| Path | What | Stack |
|------|------|-------|
| [`backend/`](backend) | Game service: runs the game loop over VK's bot API, scores answers, stores games and ratings. | Go · Gin · MongoDB |
| [`frontend/`](frontend) | Admin panel + live rating Mini App (game/topic editor, groups, leaderboard). | React · VKUI · VK Bridge |

See [`backend/README.md`](backend/README.md) for the game architecture and
deployment details.

## Quick start

Requires [mise](https://mise.jdx.dev) (pins Go + Node) and a MongoDB instance.

```bash
mise install

# backend
mise run backend:run          # game service (default :9010)

# frontend
mise run frontend:install
mise run frontend:start       # admin/rating UI dev server
```

| Task | Command |
|------|---------|
| Build backend | `mise run backend:build` |
| Test backend | `mise run backend:test` |
| Lint backend | `mise run backend:lint` |
| Build frontend | `mise run frontend:build` |
