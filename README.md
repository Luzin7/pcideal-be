# PC Ideal - Backend

> API inteligente para recomendação de builds de PC personalizadas

API REST em Go que utiliza estratégias de orçamento dinâmicas para recomendar configurações de PC otimizadas com base no perfil de uso e orçamento do usuário.

---

## Funcionalidades

- **Consulta de peças** com atualização automática de preços via scraper externo
- **Recomendações de builds** com 3 perfis (Econômica, Balanceada, Performance)
- **Estratégias de orçamento inteligentes** adaptadas ao tipo de uso (Gaming, Work, Content Creator)
- **Validação de compatibilidade** entre componentes

---

## Arquitetura

Projeto estruturado seguindo **Clean Architecture**:

- **Domain**: Entidades e regras de negócio (Part, BudgetStrategy)
- **Use Cases**: Lógica de negócio e algoritmos de seleção de componentes
- **Infrastructure**: HTTP controllers, repositórios MongoDB, integrações externas
- **Dependency Injection**: Configurada no `main.go`

---

## Como Funciona

O sistema gera 3 builds (Econômica, Balanceada, Performance) seguindo estas etapas:

1. **Seleção de estratégia de orçamento** baseada no tipo de uso e valor
2. **Busca de componentes compatíveis** seguindo ordem de dependências (CPU → Motherboard → GPU → PSU → RAM → SSD)
3. **Validação de compatibilidade** (socket, tipo de memória, potência)

---

## Endpoints

```http
GET  /health                      # Health check
GET  /api/parts/                  # Lista peças com paginação
GET  /api/parts/:id               # Detalhes de uma peça
POST /api/builds/recommendations  # Gera recomendações de builds
```

### Listagem de peças

`GET /api/parts/` aceita `limit` (padrão 50, máximo 200) e `offset` (padrão 0). Listas vazias retornam `200` com `data: []`.

```http
GET /api/parts/?limit=50&offset=0
```

**Exemplo de requisição de build:**
```json
{
  "budget_cents": 500000,
  "usage_type": "GAMING",
  "cpu_preference": "AMD",
  "gpu_preference": "NVIDIA"
}
```

---

## Stack

- **Go 1.24** com Gin Framework
- **MongoDB** para persistência
- **Scraper API** para atualização de preços

---

## Executando o Projeto

```bash
# Clone e configure
git clone https://github.com/Luzin7/pcideal-be.git
cd pcideal-be

# Configure .env com as credenciais necessárias
PORT=8080
DATABASE_URL=mongodb://localhost:27017
PCIDEAL_DB_NAME=pcideal
SCRAPER_API_URL=...
SCRAPER_API_KEY=...
# CIDRs do proxy reverso; obrigatório em produção quando há proxy
TRUSTED_PROXIES=10.0.0.0/8
ENV=development

# Inicie MongoDB
docker-compose up -d

# Execute
go run cmd/main.go
```

O servidor cria os índices do MongoDB na inicialização e encerra de forma graciosa em `SIGINT`/`SIGTERM`.

### Health check

`GET /health` faz ping no MongoDB com timeout de 2s. Retorna `503` quando o banco não responde e `200` com `status: "ok"` quando responde. O campo `unique_index_ready` indica se o índice único `parts_store_url_unique` existe; ele não derruba a readiness sozinho, mas `false` significa que a importação não é idempotente.

### TRUSTED_PROXIES

Quando a API roda atrás de um proxy reverso, configure `TRUSTED_PROXIES` com os CIDRs do proxy (separados por vírgula). Vazio faz o `ClientIP` ignorar `X-Forwarded-For`: em produção, todos os clientes atrás do proxy passam a compartilhar o mesmo bucket de rate limit. O boot registra um `WARNING` nesse caso e não falha.

### Índice único de peças (store, url)

A idempotência da importação depende do índice `parts_store_url_unique`. A criação não é fatal: se existirem duplicatas legadas, a API sobe, registra um `ERROR` com o comando de auditoria e expõe `unique_index_ready: false` no `/health`. Rode a auditoria abaixo e garanta que o resultado seja vazio **antes** de confiar na idempotência:

```js
db.parts.aggregate([
  { $group: { _id: { store: "$store", url: "$url" }, count: { $sum: 1 } } },
  { $match: { count: { $gt: 1 } } },
])
```

Enquanto a consulta retornar documentos, remova ou corrija as duplicatas antes de confiar no índice.

## Testes

```bash
go test ./...                    # Todos os testes
go test -cover ./...             # Com cobertura
go vet ./...                     # Análise estática
go build ./...                   # Compilação
```

---

**PC Ideal** - Builds inteligentes para todos os orçamentos
