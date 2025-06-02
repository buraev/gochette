# Gochette

**Gochette** — это минималистичный сервер для кэширования и проксирования данных.

### Запуск

Чтобы собрать и запустить `gochette`:

1. запоните .env

```bash
CACHE_FOLDER="data"                   # Путь до папки где будут храниться данные 
CORE_PATH="http://localhost:8000"     # URL по которому будет осуществляться доступ к api
ALLOW_FRONTEND="https://buraev.com"   # Access-Control-Allow-Origin для CORS (используйте * для всех)
GITHUB_ACCESS_TOKEN=someToken         # Токен доступа из GitHub 
VALID_TOKENS="someToken"              # Bearer token для авторизации (будет добавлен позже)
```

2. запустите команду
```bash
make start
```

---

## 📡 API

На данный момент реализовано:

- `GET /github/` — получение списка pinned-репозиториев пользователя GitHub (через GraphQL API с кэшированием).

В планах:

- 🔜 Поддержка Steam (публичные игровые профили, достижения и т.д.)
- 🔜 Интеграция с Apple Music (любимые треки, сейчас играет и т.д.)

---

## 📄 Лицензия

MIT — см. [LICENSE](./LICENSE)
