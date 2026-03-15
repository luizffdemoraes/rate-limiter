# Rate Limiter (Go)

Middleware em **Go** para controlar o fluxo de requisições de um serviço web. O sistema limita o tráfego com base no **IP do solicitante** ou em um **token de acesso**, usando **Redis** para persistência e orquestração.

---

## Objetivo

Desenvolver um **Rate Limiter** que funcione como **middleware**, capaz de:

- Limitar requisições por **endereço IP** ou por **token** (`API_KEY`).
- Persistir contagem e janelas de tempo no **Redis** (via Docker).
- Manter arquitetura extensível (padrão **Strategy**) para trocar o backend de persistência no futuro.

---

## Regras de negócio (critérios de limitação)

### Limitação por IP

Restringe o número máximo de **requisições por segundo** recebidas de um único endereço IP.

### Limitação por token

Restringe o número máximo de **requisições por segundo** com base em um token de acesso único.

- **Header:** o token é enviado como `API_KEY: <TOKEN>`.

### Precedência (regra de ouro)

As configurações do **token** prevalecem sobre as do **IP**.

**Exemplo:** limite global por IP = 10 req/s; um token específico = 100 req/s → o sistema aplica **100 req/s** para requisições com esse token.

---

## Comportamento em caso de bloqueio

Quando o limite for excedido (por IP ou por token):

| Aspecto | Comportamento |
|--------|----------------|
| **HTTP** | Status **429** imediato. |
| **Corpo** | Exatamente: `you have reached the maximum number of requests or actions allowed within a certain time frame` |
| **Tempo de bloqueio** | IP ou token infrator fica bloqueado por um período **configurável** (ex.: 5 minutos). Enquanto bloqueado, novas requisições são **rejeitadas**. |

---

## Requisitos técnicos e arquitetura

### Middleware

A lógica do rate limiter é um **middleware** que envolve o servidor HTTP.

### Persistência (Redis)

Contagem e controle de tempo ficam no **Redis** (imagem Docker no `docker-compose`).

### Design pattern: Strategy

- Implementar **Strategy** para a camada de persistência.
- **Redis é obrigatório** neste desafio, mas a arquitetura deve permitir trocar por outro storage **apenas trocando a implementação da estratégia**.

### Desacoplamento

A **regra de negócio** do limiter fica **separada** da **lógica do middleware** (ex.: middleware só orquestra; o limiter aplica regras).

### Configuração

Tudo via **variáveis de ambiente** e/ou **arquivo `.env` na raiz**, por exemplo:

- Limite máximo de requisições por segundo (IP e/ou token, conforme seu modelo).
- Tempo de bloqueio após excesso.
- Parâmetros de conexão com o Redis.

*(Documente no próprio repositório os nomes exatos das variáveis assim que implementar.)*

---

## Como executar o projeto

O avaliador deve conseguir subir aplicação e testes usando **apenas Docker / Docker Compose**.

### Subir aplicação + Redis

Na raiz do repositório (onde estão `docker-compose.yaml` e `Dockerfile`):

```bash
docker compose up --build
```

A aplicação deve escutar na porta **8080** (conforme especificação).

### Testes automatizados

Exemplo comum (ajuste ao comando real do seu `Dockerfile`/`compose`):

```bash
docker compose run --rm app go test ./...
```

Ou o target que você expuser no `docker-compose` para rodar os testes.

---

## Testes esperados

Incluir testes que demonstrem:

1. **Eficácia** do limiter (429 após ultrapassar o limite, bloqueio, etc.).
2. **Precedência token > IP** (token com limite maior que o IP deve se comportar conforme o limite do token).

---

## Como configurar o limiter

1. Copie `.env.example` para `.env` (se existir) ou defina as variáveis no `docker-compose.yaml`.
2. Ajuste limites por segundo, tempo de bloqueio e Redis conforme documentado nas envs do projeto.
3. Para **tokens com limites diferentes**, descreva no README como cadastrar/mapa token → limite (se aplicável ao seu desenho).

---

## Como alterar a estratégia de persistência

1. Defina uma **interface** (Strategy) usada pelo domínio do limiter (ex.: incrementar contador, checar bloqueio, TTL).
2. Implemente a estratégia **Redis** atual.
3. Para outro backend, crie uma nova struct que implemente a mesma interface e injete no middleware/limiter via composição ou factory — **sem alterar a regra de negócio** do limiter.

---

## Entregáveis

| Item | Descrição |
|------|-----------|
| **Código fonte** | Implementação completa do rate limiter. |
| **Dockerfile** | Build da aplicação Go. |
| **docker-compose.yaml** | Sobe a app na porta **8080** e o **Redis**. |
| **README** | Configuração, estratégias, execução (este documento + detalhes das envs após implementação). |
| **Testes** | Cobrindo limiter e precedência token > IP. |

---

## Regras de entrega

- **Repositório único:** apenas este projeto.
- **Branch principal:** todo o código na branch **`main`**.
- **Execução:** rodar projeto e testes com **Docker / Docker Compose** apenas.

---

## Licença

Defina a licença conforme a política do curso ou da sua organização.
