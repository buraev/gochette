
# Gochette

**Gochette** is a minimalist server for caching and proxying data.

### Getting Started

To build and run `gochette`:

1. Fill out the `.env` file:

```bash
CACHE_FOLDER="data"                   # Path to the folder where cached data will be stored
CORE_PATH="http://localhost:8000"     # URL where the API will be accessed
ALLOW_FRONTEND="https://buraev.com"   # Access-Control-Allow-Origin for CORS (use * for all origins)
GITHUB_ACCESS_TOKEN=someToken         # GitHub access token
VALID_TOKENS="someToken"              # Bearer token for authorization (to be implemented)
```

2. Run the command:

```bash
make start
```

---

## 📡 API

Currently implemented:

- `GET /github/` — fetches the list of pinned GitHub repositories (via GitHub GraphQL API with caching).

Planned features:

- 🔜 Steam support (public profiles, achievements, etc.)
- 🔜 Apple Music integration (favorite tracks, currently playing, etc.)

---

## 📄 License

MIT — see [LICENSE](./LICENSE)



