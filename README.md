## Rate Limiter (Go)

Middleware em **Go** para controlar o fluxo de requisições de um serviço web.  
O sistema limita o tráfego com base no **IP do solicitante** ou em um **token de acesso**, usando **Redis** para persistência e orquestração.

> O enunciado completo do desafio está em `README_CHALLENGE.md`.  
> Este `README.md` descreve **a implementação deste repositório**: como ela funciona, como configurar e como executar.

---

### Arquitetura da solução

- **`config`**: leitura das variáveis de ambiente (`.env`) em uma struct `Config`.
- **`internal/limiter`**:
  - `store.go`: interface `Store` (Strategy) para persistência.
  - `redis_store.go`: implementação `RedisStore` usando Redis.
  - `limiter.go`: regra de negócio do rate limit (IP, token, precedência, bloqueio).
- **`internal/middleware`**:
  - `ratelimit.go`: middleware HTTP (`net/http`) que extrai IP/token, chama o limiter e decide entre 200/429.
- **`internal/handler`**:
  - `handler.go`: rotas simples para exercício (`/healthz`, `/api/v1/example`).
- **`cmd/server`**:
  - `main.go`: fio da aplicação → carrega config, cria RedisStore, instancia Limiter, aplica middleware e sobe o servidor HTTP.

Pontos importantes:

- A **regra de negócio** do limiter (IP/token, precedência, bloqueio) está isolada no pacote `internal/limiter`.
- O **middleware** em `internal/middleware` não conhece detalhes de Redis; ele apenas orquestra a chamada ao limiter.
- A **Strategy** de persistência é abstraída pela interface `Store`, permitindo trocar Redis por outro backend no futuro.

---

### Configuração via variáveis de ambiente

As configurações são lidas pela função `config.Load()` a partir das envs (e do `.env`, se presente).  
Há um arquivo `.env.example` na raiz do projeto com os valores padrão.

**Variáveis de ambiente usadas:**

| Variável        | Descrição                                                   | Default       |
|-----------------|-------------------------------------------------------------|---------------|
| `HTTP_PORT`     | Porta HTTP da aplicação                                     | `8080`        |
| `RATE_LIMIT_IP` | Limite de requisições por segundo por **IP**               | `10`          |
| `RATE_LIMIT_TOKEN` | Limite de requisições por segundo por **token** (`API_KEY`) | `100`     |
| `BLOCK_DURATION`| Duração do bloqueio em **segundos** após ultrapassar limite | `300` (5 min) |
| `REDIS_HOST`    | Host do Redis                                               | `redis`       |
| `REDIS_PORT`    | Porta do Redis                                              | `6379`        |
| `REDIS_DB`      | Banco (DB) numérico do Redis                                | `0`           |
| `REDIS_PASSWORD`| Senha do Redis (se houver)                                  | vazio         |

**Passo a passo:**

1. Copie o arquivo de exemplo:

   ```bash
   cp .env.example .env
   ```

2. Ajuste os valores conforme necessidade (limites, duração de bloqueio, host/porta do Redis).

---

### Regras de negócio implementadas

- **Limitação por IP**:
  - Para requisições **sem header `API_KEY`**, o limiter utiliza o limite configurado em `RATE_LIMIT_IP`.
  - A contagem é por janela de **1 segundo** por IP.

- **Limitação por token**:
  - Para requisições com header `API_KEY: <TOKEN>`, o limiter utiliza o limite configurado em `RATE_LIMIT_TOKEN`.
  - A contagem é por janela de **1 segundo** por token.

- **Precedência token > IP**:
  - Se um token estiver presente, **sempre** prevalece o limite do token (`RATE_LIMIT_TOKEN`), independentemente do limite para o IP de origem.

- **Persistência em Redis**:
  - Contadores de requisições (por IP ou token) e flags de bloqueio são armazenados no Redis via `RedisStore` (Strategy `Store`).
  - Chaves típicas:
    - Contagem: `rate:ip:<ip>` ou `rate:token:<token>`.
    - Bloqueio: `block:ip:<ip>` ou `block:token:<token>`.

- **Bloqueio e resposta HTTP**:
  - Quando o limite é ultrapassado, o limiter grava uma chave de bloqueio com TTL = `BLOCK_DURATION`.
  - Enquanto o IP/token estiver bloqueado, novas requisições são negadas imediatamente.
  - O middleware responde com:
    - **HTTP 429 (Too Many Requests)**.
    - Corpo **exato**:  
      `you have reached the maximum number of requests or actions allowed within a certain time frame`

---

### Como executar o projeto (Docker / Docker Compose)

> Pré-requisitos: Docker e Docker Compose instalados.

Na raiz do repositório (onde estão `docker-compose.yaml` e `Dockerfile`):

```bash
docker compose up --build
```

- Isso irá:
  - Subir um container `redis`.
  - Buildar a imagem da aplicação Go.
  - Subir o container `app` escutando na porta **8080**.

**Endpoints:**

- `GET /healthz` → checagem de saúde simples.
- `GET /api/v1/example` → endpoint protegido pelo middleware de rate limit.

Exemplo rápido com `curl`:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/example
```

---

### Testes automatizados

Os testes podem ser executados diretamente com Go ou via Docker Compose.

**Go (local):**

```bash
go test ./...
```

**Via Docker Compose:**

```bash
docker compose run --rm app go test ./...
```

Os testes cobrem:

- `internal/limiter`:
  - Limite por IP e bloqueio.
  - Precedência token > IP.
- `internal/middleware`:
  - Resposta HTTP 429 com corpo exato quando o limite é estourado.
- `internal/handler`:
  - Comportamento dos handlers (`/healthz`, `/api/v1/example`).
- `config`:
  - Uso de valores padrão e leitura correta das envs.

---

### Strategy de persistência (trocando Redis por outro backend)

- A interface `Store` em `internal/limiter/store.go` define as operações necessárias:
  - `Increment(key string, windowSeconds int) (int64, error)`
  - `IsBlocked(key string) (bool, error)`
  - `Block(key string, duration time.Duration) error`
- A implementação atual (`RedisStore`) em `internal/limiter/redis_store.go` usa Redis.
- Para trocar o backend de persistência:
  1. Crie uma nova struct que implemente `Store` (por ex., `MemoryStore`, `SQLStore`).
  2. No `main.go`, substitua a criação de `RedisStore` pela nova implementação.
  3. Não é necessário alterar o `Limiter` nem o middleware.

---

### Testes manuais (curl / Postman)

Além dos testes automatizados, você pode validar o comportamento via:

- **Collection Postman**:
  - Arquivo: `collection/rate-limiter.postman_collection.json`.
  - Contém requisições para `/healthz`, `/api/v1/example` (com e sem `API_KEY`).

- **Exemplo de cenário com `curl` (limitação por IP)**:

  ```bash
  for i in {1..20}; do
    echo "Request $i:"
    curl -is -w "HTTP %{http_code}\n" -o /dev/null "http://localhost:8080/api/v1/example"
  done
  ```

- **Exemplo de cenário com `curl` (limitação por token)**:

  ```bash
  for i in {1..20}; do
    echo "Request $i:"
    curl -is -w "HTTP %{http_code}\n" -o /dev/null \
      -H "API_KEY: my-token" \
      "http://localhost:8080/api/v1/example"
  done
  ```

