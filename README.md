# CLI de teste de carga (Go)

Sistema de **CLI (Command Line Interface)** em Go para realizar **testes de carga** em serviços web. O usuário informa a URL do serviço, o número total de requisições e a quantidade de chamadas simultâneas. Ao final, a aplicação gera um **relatório detalhado** da execução no console.

---

## Objetivo do desafio

- Executar requisições HTTP contra um serviço alvo.
- Respeitar o total de requisições e o nível de concorrência definidos.
- Exibir métricas claras ao término do teste.

---

## Entrada de parâmetros

A aplicação deve aceitar os seguintes parâmetros via linha de comando:

| Parâmetro       | Descrição                                      |
|-----------------|------------------------------------------------|
| `--url`         | URL do serviço a ser testado.                  |
| `--requests`    | Número total de requisições a serem realizadas.|
| `--concurrency` | Número de chamadas simultâneas.                |

---

## Requisitos técnicos

### 1. Execução do teste

- Realizar requisições **HTTP** para a URL especificada.
- **Distribuir** as requisições de acordo com o nível de **concorrência** definido.
- Garantir que o número total de requisições (`--requests`) seja cumprido **exatamente**.

### 2. Relatório (console)

Ao final da execução, o sistema deve apresentar:

- **Tempo total** gasto na execução.
- **Quantidade total** de requests realizados.
- **Quantidade** de requests com status HTTP **200**.
- **Distribuição** dos demais códigos de status HTTP (ex.: quantidade de 404, 500, etc.).

---

## Execução via Docker (obrigatório)

A aplicação deve ser empacotada em uma **imagem Docker** para facilitar a execução.

Exemplo de comando de teste:

```bash
docker run <sua-imagem-docker> --url=http://google.com --requests=1000 --concurrency=10
```

Substitua `<sua-imagem-docker>` pelo nome/tag da imagem que você construiu (ex.: `loadtest-cli:latest`).

---

## Build da imagem Docker

Na raiz do repositório (onde está o `Dockerfile`):

```bash
docker build -t loadtest-cli:latest .
```

*(Ajuste o nome da tag `loadtest-cli:latest` conforme preferir.)*

---

## Executar localmente (Go)

Se quiser rodar sem Docker (após implementar o CLI):

```bash
go run . --url=https://exemplo.com --requests=100 --concurrency=5
```

Ou, após compilar:

```bash
go build -o loadtest .
./loadtest --url=https://exemplo.com --requests=100 --concurrency=5
```

---

## Entregáveis

| Item        | Descrição                                                |
|------------|-----------------------------------------------------------|
| Código fonte | Repositório com a implementação do CLI.                |
| `Dockerfile` | Configuração para construção da imagem.                  |
| `README`   | Instruções de como buildar a imagem e executar o teste.  |

---

## Regras de entrega

- **Repositório exclusivo:** o repositório deve conter **apenas** este projeto.
- **Branch principal:** todo o código deve estar na branch **`main`**.

---

## Licença

Defina a licença do projeto conforme a política do curso ou da sua organização.
