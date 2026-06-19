# VK Jeopardy

Backend for a **Jeopardy-style quiz game** that runs inside [VKontakte](https://vk.com)
community chats. Players pick a topic and point value, the bot posts the question, and
the first correct answer scores. The service drives the game loop over VK's Bots Long
Poll / Callback API, persists games and scores in MongoDB, and serves a live rating
page (a VK Mini App) for each community.

## How it works

```
VK community  ──callback──▶  Gin HTTP API ──▶  game engine ──▶  MongoDB
   message                    (/api/...)        (answer match,    (games,
                                  │              scoring, retry)    users, scores)
                                  ▼
                          VK API replier  ──▶  posts answers back to the chat
                          (rate-limited, flood-control retry)
```

- **`game/`** — the game engine: turn handling, answer matching, retry/branch rules.
- **`domain/`** — entities, MongoDB DAO, the VK API client and the rate-limited
  `replier` (batches replies and retries on VK flood-control errors).
- **`rating/`** — read API for the per-community rating Mini App, with HMAC-signed
  VK launch-param auth and in-memory rating/name caches.
- **`configuration/`** — strict YAML config loader (`application.yml`).
- **`ansible/`** — provisioning + zero-downtime deploy playbooks (nginx + systemd,
  Let's Encrypt).

A game is authored as JSON (see `docs/games.json.js`): topics, point tiers, questions
with multiple accepted answers, custom messages, and rules (`numTries`, `instantWin`).

## Tech stack

Go 1.23 · [Gin](https://github.com/gin-gonic/gin) · MongoDB
([mongo-driver](https://go.mongodb.org/mongo-driver)) · VK API · Ansible.

## Local development

Requires [mise](https://mise.jdx.dev) (pins the Go toolchain) and a MongoDB instance.

```bash
mise install                 # install the pinned Go toolchain
cp application.example.yml application.yml   # then edit Mongo/VK credentials
mise run run                 # start the server (default :9010)
```

Common tasks:

| Task | Command |
|------|---------|
| Build binary | `mise run build` |
| Run tests | `mise run test` |
| Lint (vet + golangci-lint) | `mise run lint` |
| Tidy modules | `mise run tidy` |

`mockResponse: true` in the config lets the engine run without posting back to VK —
handy for local testing.

## Deployment

Provision a fresh Ubuntu host, then deploy with the bundled Ansible playbooks:

```bash
cd ansible
cp hosts_example.yml hosts/prod.yml          # fill in your inventory
ansible-playbook playbook.yml -i hosts/prod.yml          # one-time provisioning

# build + ship the backend
GOOS=linux GOARCH=amd64 go build -o deploy/goj .
./deploy.sh b                                # backend   (deploy_playbook.yml)
./deploy.sh f                                # frontend  (deploy_frn_playbook.yml)
```
