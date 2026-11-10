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
# CIDRs do proxy reverso; alternativa ao CLIENT_IP_HEADER
TRUSTED_PROXIES=10.0.0.0/8
# Header do IP real do cliente; recomendado atrás do Cloudflare
CLIENT_IP_HEADER=
# Origens de CORS; vazio usa https://pcideal-web.vercel.app
CORS_ALLOWED_ORIGINS=
ENV=development

# Inicie MongoDB
docker-compose up -d

# Execute
go run cmd/main.go
```

O servidor cria os índices do MongoDB na inicialização e encerra de forma graciosa em `SIGINT`/`SIGTERM`.

### Health check

`GET /health` faz ping no MongoDB com timeout de 2s. Retorna `503` quando o banco não responde e `200` com `status: "ok"` quando responde. O campo `unique_index_ready` indica se o índice único `parts_store_url_unique` existe; ele não derruba a readiness sozinho, mas `false` significa que a importação não é idempotente.

### Proxy, IP do cliente e CORS

Topologia de produção: `Cloudflare → Render`. A API pública é `https://pcideal-be.onrender.com` e o frontend publicado é `https://pcideal-web.vercel.app`. O Render não publica CIDRs dos proxies de entrada, então não existe lista de CIDRs do Render para `TRUSTED_PROXIES`; a única forma de confiar no hop do Render é confiar em toda a internet.

`X-Forwarded-For` só é confiável quando o hop que o escreve é conhecido. Duas opções:

**Opção A — `TRUSTED_PROXIES`.** Só use se a origem não for alcançável diretamente: `TRUSTED_PROXIES=0.0.0.0/0,::/0`. Hoje a origem responde `403` para acesso direto. Se esse bloqueio cair, qualquer cliente forja `X-Forwarded-For` e escapa do rate limit por IP.

**Opção B — `CLIENT_IP_HEADER` (recomendada).** `CLIENT_IP_HEADER=CF-Connecting-IP`. O Cloudflare sobrescreve esse header na borda, então o cliente não consegue forjá-lo. A API valida o valor com `net.ParseIP`; header ausente ou inválido cai para o `ClientIP()` do Gin, e um valor com vírgula é recusado no boot. Habilite apenas atrás de uma borda que de fato sobrescreve o header.

Os ranges oficiais do Cloudflare (`https://www.cloudflare.com/ips-v4` e `https://www.cloudflare.com/ips-v6`) servem para allowlist de firewall, não como `TRUSTED_PROXIES`: são os IPs de saída do Cloudflare, não os hops internos até o Render.

**CORS.** `CORS_ALLOWED_ORIGINS` lista as origens permitidas separadas por vírgula, como origens absolutas (`esquema://host`). Espaços e barra final são normalizados. Vazio usa o padrão `https://pcideal-web.vercel.app`, a única origem de produção. Desenvolvimento local pode definir `CORS_ALLOWED_ORIGINS=http://localhost:3000` no `.env` não versionado, nunca em produção. Qualquer outra origem recebe `403`, inclusive `localhost` e `pcideal.com.br`; a permitida é ecoada em `Access-Control-Allow-Origin`. Configuração inválida derruba o boot.

### Índice único de peças (store, url)

A idempotência da importação depende do índice `parts_store_url_unique`. A criação não é fatal: se existirem duplicatas legadas, a API sobe, registra um `ERROR` com o comando de auditoria e expõe `unique_index_ready: false` no `/health`. Rode a auditoria abaixo e garanta que o resultado seja vazio **antes** de confiar na idempotência:

```js
db.parts.aggregate([
  { $group: { _id: { store: "$store", url: "$url" }, count: { $sum: 1 } } },
  { $match: { count: { $gt: 1 } } },
])
```

Enquanto a consulta retornar documentos, remova ou corrija as duplicatas antes de confiar no índice.

O comando `cmd/dedupe-parts` audita e corrige as duplicatas legadas. O modo padrão é dry-run e não escreve nada:

```bash
go run ./cmd/dedupe-parts                 # dry-run: grupos, exclusões, campos divergentes e store/url ausentes
go run ./cmd/dedupe-parts --check         # dry-run: sai com código 1 quando existe grupo duplicado
go run ./cmd/dedupe-parts --timeout=10m   # limite de execução (padrão: 5m)
go run ./cmd/dedupe-parts --apply --confirm=DELETE-DUPLICATES
```

`--apply` exige `--confirm=DELETE-DUPLICATES` e recusa a execução quando há documentos com `store` ou `url` vazio/ausente. Em cada grupo, mantém o documento com maior `updated_at` (empate: maior `_id`) e remove somente os `_id` extras, com filtro também por `store` e `url`; depois cria ou valida o índice `parts_store_url_unique` e confere a retenção dos canônicos e a contagem de exclusões. Se a execução for interrompida, rode novamente: a operação é idempotente e conclui o restante.

A importação deve ser idempotente e usar `mongoimport --upsert --upsertFields store,url`, para que o índice único seja respeitado.

## Testes

```bash
go test ./...                    # Todos os testes
go test -cover ./...             # Com cobertura
go vet ./...                     # Análise estática
go build ./...                   # Compilação
```

---

**PC Ideal** - Builds inteligentes para todos os orçamentos
